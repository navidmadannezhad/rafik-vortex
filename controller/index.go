package controller

import (
	"airun/code-reviewer/models"
	"airun/code-reviewer/services"
	"airun/code-reviewer/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequestPRReview(requestContext *gin.Context) {
	body := models.PRWebhookPayload{}
	headers := map[string]interface{}{}

	err := requestContext.BindJSON(&body)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	err2 := requestContext.BindHeader(&headers)
	if err2 != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err2.Error(),
			},
		)
		return
	}

	isPullRequestEvent := body.PullRequest != nil
	if !isPullRequestEvent {
		requestContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ResolveErrorMessage("Not a pull request event"),
		)
		return
	}

	// pr = payload["pull_request"]
	// repo = payload["repository"]
	// owner = repo["owner"]["login"]
	// repo_name = repo["name"]
	// pr_number = pr["number"]
	// head_sha = pr["head"]["sha"]

	pullRequest := body.PullRequest
	repository := body.Repository
	owner := repository.Owner
	repositoryName := repository.Name
	pullRequestNumber := pullRequest.Number
	// headSha := pullRequest.Head.SHA
	fmt.Printf("Received PR! Number: %d, Repository: %s\n", pullRequestNumber, repositoryName)

	resp, err := services.GetPullRequestDiffQuery(
		models.GitDiffApiUrlOptions{
			PullRequestNumber: pullRequestNumber,
			RepositoryName:    repository.Name,
			OwnerName:         owner.Login,
		},
	)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}

	rawResponseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}

	files := utils.ParsePullRequestDiff(string(rawResponseBody))

	var aiResponseBody = models.GetModelReviewMutationResponse{}
	payload := models.GetModelReviewMutationOptions{
		Files: files,
	}
	aiResponse, err := services.GetModelReviewMutation(payload)
	if err != nil || aiResponse.StatusCode != 200 {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}
	defer aiResponse.Body.Close()

	err = json.NewDecoder(aiResponse.Body).Decode(&aiResponseBody)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}
	comments := utils.ParseAIReviewComments(aiResponseBody.Choices[0].Message.Content)

	postReviewResponse, err := services.PostReviewToGitMutation(
		models.PostReviewToGitMutationOptions{
			Owner:             owner.Login,
			Repository:        repositoryName,
			PullRequestNumber: pullRequestNumber,
			HeadSha:           pullRequest.Head.SHA,
			Comments:          comments,
		},
	)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}
	defer postReviewResponse.Body.Close()

	requestContext.JSON(
		http.StatusOK,
		utils.ResolveSuccessMessage("Review completed"),
	)
}

func getRuleFiles(stack models.RepoTechnology) ([]LibraryFile, error) {
	rootPath := fmt.Sprintf("/%v/rules", stack)
	getDirectoryPayload := models.GetProjectDirectoryListQueryOptions{
		OwnerName:      utils.GetEnv("LIBRARY_REPO_OWNER", "test-org"),
		RepositoryName: utils.GetEnv("LIBRARY_REPO_NAME", "library"),
		RootPath:       &rootPath,
	}
	ruleFilesResponse, err := services.GetProjectDirectoryListQuery(getDirectoryPayload)
	if err != nil {
		return nil, err
	}
	defer ruleFilesResponse.Body.Close()
	var ruleFiles models.GetProjectDirectoryListQueryResponse
	err = json.NewDecoder(ruleFilesResponse.Body).Decode(&ruleFiles)
	if err != nil {
		return nil, err
	}

	libraryRuleFiles := []LibraryFile{}
	for _, file := range ruleFiles.DirContents {
		fileContentPayload := models.GetGitFileContentQueryOptions{
			OwnerName:      utils.GetEnv("LIBRARY_REPO_OWNER", "test-org"),
			RepositoryName: utils.GetEnv("LIBRARY_REPO_NAME", "library"),
			FilePath:       file.Path,
		}
		fileContentResponse, err := services.GetGitFileContentQuery(fileContentPayload)
		if err != nil {
			return nil, err
		}
		defer fileContentResponse.Body.Close()
		bodyBytes, err := io.ReadAll(fileContentResponse.Body)
		if err != nil {
			return nil, err
		}

		libraryRuleFile := LibraryFile{
			Content:  string(bodyBytes),
			FilePath: "rules/" + file.Name,
		}

		libraryRuleFiles = append(libraryRuleFiles, libraryRuleFile)
	}

	return libraryRuleFiles, nil
}

func getSkillFiles(stack models.RepoTechnology) ([]LibraryFile, error) {
	rootPath := fmt.Sprintf("/%v/skills", stack)
	getDirectoryPayload := models.GetProjectDirectoryListQueryOptions{
		OwnerName:      utils.GetEnv("LIBRARY_REPO_OWNER", "test-org"),
		RepositoryName: utils.GetEnv("LIBRARY_REPO_NAME", "library"),
		RootPath:       &rootPath,
	}
	skillFoldersResponse, err := services.GetProjectDirectoryListQuery(getDirectoryPayload)
	if err != nil {
		return nil, err
	}
	defer skillFoldersResponse.Body.Close()
	var skillFolders models.GetProjectDirectoryListQueryResponse
	err = json.NewDecoder(skillFoldersResponse.Body).Decode(&skillFolders)
	if err != nil {
		return nil, err
	}

	librarySkillFiles := []LibraryFile{}
	for _, skillFolder := range skillFolders.DirContents {
		getSkillContentPayload := models.GetGitFileContentQueryOptions{
			OwnerName:      utils.GetEnv("LIBRARY_REPO_OWNER", "test-org"),
			RepositoryName: utils.GetEnv("LIBRARY_REPO_NAME", "library"),
			FilePath:       string(stack) + "/skills/" + skillFolder.Name + "/SKILL.md",
		}
		getSkillContentResponse, err := services.GetGitFileContentQuery(getSkillContentPayload)
		if err != nil {
			return nil, err
		}
		defer getSkillContentResponse.Body.Close()

		bodyBytes, err := io.ReadAll(getSkillContentResponse.Body)
		if err != nil {
			return nil, err
		}
		librarySkillFiles = append(librarySkillFiles, LibraryFile{
			Content:  string(bodyBytes),
			FilePath: "skills/" + skillFolder.Name + "/SKILL.md",
		})
	}

	return librarySkillFiles, nil
}

type LibraryFile struct {
	Content  string `json:"content"`
	FilePath string `json:"file_path"`
}

func getRelatedInstructionFilesFromLibrary(stack models.RepoTechnology) ([]LibraryFile, error) {
	relatedFiles := []LibraryFile{}
	libraryRootOptions := models.GetProjectDirectoryListQueryOptions{
		OwnerName:      utils.GetEnv("LIBRARY_REPO_OWNER", "test-org"),
		RepositoryName: utils.GetEnv("LIBRARY_REPO_NAME", "library"),
	}
	libraryRootResponse, err := services.GetProjectDirectoryListQuery(libraryRootOptions)
	if err != nil {
		return []LibraryFile{}, err
	}
	var rootDirsResponse = models.GetProjectDirectoryListQueryResponse{}
	err = json.NewDecoder(libraryRootResponse.Body).Decode(&rootDirsResponse)
	if err != nil {
		log.Print("ooor hereee?")
		return []LibraryFile{}, err
	}

	stackExistsInLibrary := false
	for _, rootDir := range rootDirsResponse.DirContents {
		if rootDir.Name == string(stack) {
			stackExistsInLibrary = true
		}
	}
	if !stackExistsInLibrary {
		return []LibraryFile{}, errors.New("Stack does not exist in library.")
	}

	ruleFiles, err := getRuleFiles(stack)
	if err != nil {
		log.Print("HERE")
		return nil, err
	}
	skillFiles, err := getSkillFiles(stack)
	if err != nil {
		log.Print("HERE2")
		return nil, err
	}

	relatedFiles = append(relatedFiles, ruleFiles...)
	relatedFiles = append(relatedFiles, skillFiles...)

	return relatedFiles, nil
}

func libraryFilesToChangeFileOperations(libFiles []LibraryFile) []models.ChangeFileOperation {
	var fileOperations []models.ChangeFileOperation
	for _, libFile := range libFiles {
		fileOperations = append(fileOperations, models.ChangeFileOperation{
			Content:   libFile.Content,
			Operation: "create",
			Path:      libFile.FilePath,
		})
	}

	return fileOperations
}

func getFileOperations(targetRepoName string, targetRepoOwner string, relatedFiles []LibraryFile) ([]models.ChangeFileOperation, error) {
	fileOperations := []models.ChangeFileOperation{}
	for _, relatedFile := range relatedFiles {
		var relatedFilePath string
		if strings.Contains(relatedFile.FilePath, "skills") {
			relatedFilePath = ".cursor/" + relatedFile.FilePath
		} else if strings.Contains(relatedFile.FilePath, "rules") {
			relatedFilePath = ".cursor/" + relatedFile.FilePath
		}
		targetRepoFileMetadataPayload := models.GetProjectDirectoryListQueryOptions{
			OwnerName:      targetRepoOwner,
			RepositoryName: targetRepoName,
			RootPath:       &relatedFilePath,
		}
		targetRepoFileMetadataResponse, err := services.GetProjectDirectoryListQuery(targetRepoFileMetadataPayload)
		if err != nil {
			return nil, err
		}
		defer targetRepoFileMetadataResponse.Body.Close()

		if targetRepoFileMetadataResponse.StatusCode == 404 {
			// file is newly getting added to repo
			fileOperations = append(fileOperations, models.ChangeFileOperation{
				Content:   utils.GetBased64From(relatedFile.Content),
				Operation: models.Create,
				Path:      ".cursor/" + relatedFile.FilePath,
			})
		} else if targetRepoFileMetadataResponse.StatusCode == 200 {
			// file already exists in repo
			targetRepoFileMetadata := models.GetProjectDirectoryListQueryResponse{}
			err = json.NewDecoder(targetRepoFileMetadataResponse.Body).Decode(&targetRepoFileMetadata)
			if err != nil {
				return nil, err
			}
			fileOperations = append(fileOperations, models.ChangeFileOperation{
				Content:   utils.GetBased64From(relatedFile.Content),
				Operation: models.Update,
				Path:      ".cursor/" + relatedFile.FilePath,
				Sha:       targetRepoFileMetadata.FileContents.Sha,
			})
		}
	}

	return fileOperations, nil
}

func RequestAIInstructions(requestContext *gin.Context) {
	body := models.RepoPushWebhookPayload{}
	err := requestContext.BindJSON(&body)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}

	isAIUpdateInstructPush := utils.ExistsUpdateAIArg(body)
	if !isAIUpdateInstructPush {
		return
	}

	log.Printf("Received Push Event for %v/%v. Checking for related instruction files ...", body.Repository.Owner.Login, body.Repository.Name)
	projectRootList, err := services.GetProjectDirectoryListQuery(
		models.GetProjectDirectoryListQueryOptions{
			OwnerName:      body.Repository.Owner.Login,
			RepositoryName: body.Repository.Name,
		},
	)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}

	var gitFilesResponse = models.GetProjectDirectoryListQueryResponse{}
	err = json.NewDecoder(projectRootList.Body).Decode(&gitFilesResponse)
	if err != nil {
		log.Print("problem here")
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}

	log.Printf("Getting technology detection file for %v/%v ...", body.Repository.Owner.Login, body.Repository.Name)
	technologyDetectedFile := utils.FindTechnologyDetectionFile(gitFilesResponse.DirContents)
	if technologyDetectedFile == "" {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage("Technology detection file not found"),
		)
		return
	}

	log.Printf("Getting Tech file content ...")
	techFileContentResponse, err := services.GetGitFileContentQuery(
		models.GetGitFileContentQueryOptions{
			FilePath:       technologyDetectedFile,
			OwnerName:      body.Repository.Owner.Login,
			RepositoryName: body.Repository.Name,
		},
	)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}
	defer techFileContentResponse.Body.Close()

	rawTechFileContent, err := io.ReadAll(techFileContentResponse.Body)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}

	techFileContent := string(rawTechFileContent)
	repoStack := utils.FindRepoStack(techFileContent)

	log.Printf("Repo stack detected: %v. Getting related instruction files from library ...", repoStack)
	relatedFiles, err := getRelatedInstructionFilesFromLibrary(repoStack)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}

	log.Printf("Found related instruction files. Getting the files ...")
	changeFileOperations, err := getFileOperations(
		body.Repository.Name,
		body.Repository.Owner.Login,
		relatedFiles,
	)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}

	log.Printf("Files recieved. Preparing to modify repository ...")
	modifyFileInRepoPayload := models.ModifyFilesInRepoMutationOptions{
		OwnerName:      body.Repository.Owner.Login,
		RepositoryName: body.Repository.Name,
		Body: models.ModifyFilesInRepoBody{
			Files:   changeFileOperations,
			Message: "feat: add some ai instructions for you comrade!",
			Branch:  utils.ExtractBranchName(body.Ref),
		},
	}

	_, err = services.ModifyFilesInRepoMutation(modifyFileInRepoPayload)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}

	log.Printf("AI instructions added successfully.")
	requestContext.JSON(
		http.StatusOK,
		utils.ResolveSuccessMessage("Instructions requested"),
	)
}
