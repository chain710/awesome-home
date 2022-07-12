package cmd

import (
	"context"
	"errors"
	"github.com/chain710/awesome-home/cmd/values"
	"github.com/chain710/awesome-home/internal/log"
	"github.com/chain710/awesome-home/internal/ver"
	dockertypes "github.com/docker/docker/api/types"
	dockercli "github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

type peekLatestImageCmd struct {
	timeout   time.Duration
	tagRegexp values.Regexp
}

func (c *peekLatestImageCmd) findLatestImageTag(local name.Tag) (string, error) {
	tagStr := local.TagStr()
	if tagStr == name.DefaultTag {
		// if local use latest, return latest
		return tagStr, nil
	}

	localVersion, err := ver.NewGeneric(tagStr)
	if err != nil {
		return "", err
	}

	var versions ver.GenericVersions
	if tags, err := remote.List(local.Repository); err != nil {
		log.Errorf("list remote image %s error: %s", local.Repository.Name(), err)
		return "", err
	} else {
		for _, tagName := range tags {
			if !c.tagRegexp.MatchString(tagName) {
				log.Debugf("tag %s not match filter", tagName)
				continue
			}
			version, err := ver.NewGeneric(tagName)
			if err != nil {
				log.Errorf("new version from tag %s error %s", tagName, err)
				return "", err
			}

			versions = append(versions, *version)
		}
	}

	if len(versions) == 0 {
		return "", errors.New("no available versions")
	}
	sort.Sort(versions)
	find := versions.UpperBound(localVersion)
	return find.String(), nil
}

func (c *peekLatestImageCmd) RunE(cmd *cobra.Command, _ []string) error {
	dockerClient, err := dockercli.NewClientWithOpts()
	if err != nil {
		log.Errorf("create docker client error: %s", err)
		return err
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), c.timeout)
	defer cancel()

	listOption := dockertypes.ContainerListOptions{All: true}
	containers, err := dockerClient.ContainerList(ctx, listOption)
	if err != nil {
		log.Errorf("list container error: %s", err)
		return err
	}

	for _, container := range containers {
		log.Debugf("image=%s id=%s", container.Image, container.ImageID)
		tag, err := name.NewTag(container.Image)
		if err != nil {
			log.Errorf("invalid image tag: %s, err: %s", container.Image, err)
			return err
		}

		latestTag, err := c.findLatestImageTag(tag)
		if err != nil {
			return err
		}
		cmd.Printf("Container %s has newer tag \"%s\"\n", container.Image, latestTag)
	}

	return nil
}

func init() {
	cmd := peekLatestImageCmd{
		timeout: time.Minute,
	}
	realCmd := &cobra.Command{
		Use:  "peek_latest_image",
		RunE: cmd.RunE,
	}
	rootCmd.AddCommand(realCmd)
	realCmd.Flags().DurationVar(&cmd.timeout, "timeout", cmd.timeout, "timeout")
	realCmd.Flags().VarP(&cmd.tagRegexp, "tag-filter", "f", "regexp for registry tag filtering")
}
