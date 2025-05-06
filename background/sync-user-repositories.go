package background

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/kabaka"
	"github.com/weichen-lin/stargazer/infrastructure"
	"github.com/weichen-lin/stargazer/util"
)

type SyncUserRepositoriesMessage struct {
	ClerkId string `json:"clerk_id"`
	Page    int    `json:"page"`
}

func (b *Background) PublishSyncUserRepositoriesEvent(ctx *gin.Context) {
	clerkId, ok := infrastructure.GetClerkId(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "not found clerk id at context"})
		return
	}
	time.Sleep(time.Second * 3)

	if canSync := b.service.GetSyncLimiter(ctx, clerkId); !canSync {
		ctx.JSON(http.StatusConflict, gin.H{"error": "sync is in progress"})
		return
	}

	page := 1

	req := &SyncUserRepositoriesMessage{
		ClerkId: clerkId,
		Page:    page,
	}

	msg, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := b.kabaka.Publish(string(SyncUserRepositories), msg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (b *Background) SyncUserRepositories(msg *kabaka.Message) error {

	req := &SyncUserRepositoriesMessage{}
	if err := json.Unmarshal(msg.Value, req); err != nil {
		return err
	}

	ctx := context.WithValue(context.Background(), "clerkId", req.ClerkId)

	token, err := b.service.GetOauthToken(ctx)
	if err != nil {
		return err
	}

	repositories, err := util.GetUserStarredRepos(req.Page, token)
	if err != nil {
		return err
	}

	for _, repository := range repositories {
		err := b.service.SaveGithubRepository(ctx, &repository)
		if err != nil {
			return err
		}

		err = b.service.SaveUserStarredRepository(ctx, int32(repository.ID))
		if err != nil {
			return err
		}
	}

	if len(repositories) >= 30 {
		req.Page++

		jsonString, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("error marshalling JSON: %s", err.Error())
		}

		b.kabaka.Publish(string(SyncUserRepositories), jsonString)
		return nil

	}

	starsCount := (req.Page-1)*30 + len(repositories)

	info, err := b.service.GetUserInfo(ctx)
	util.SendMail(&util.SendMailParams{
		Email:      info.Email,
		Name:       info.Name,
		StarsCount: starsCount,
	}, b.config.ResendApiKey)

	return b.service.UpdateUserCrontab(ctx, int32(starsCount))
}
