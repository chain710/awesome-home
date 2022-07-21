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
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type peekLatestImageCmd struct {
	timeout   time.Duration
	tagRegexp values.Regexp
}

type imageTag struct {
	name   string
	digest string
}

func (c *peekLatestImageCmd) getLatestImageTag(local name.Tag) (imageTag, error) {
	tagName, err := c.getLatestImageTagName(local)
	if err != nil {
		return imageTag{}, err
	}

	if digest, err := c.getRemoteDigest(local.Tag(tagName)); err != nil {
		return imageTag{}, err
	} else {
		return imageTag{name: tagName, digest: digest}, nil
	}
}

func (c *peekLatestImageCmd) getLatestImageTagName(local name.Tag) (string, error) {
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

func (c *peekLatestImageCmd) getRemoteDigest(ref name.Reference) (string, error) {
	if imageSpec, err := remote.Image(ref); err != nil {
		log.Errorf("inspect remote image %s error %s", ref.Name(), err)
		return "", err
	} else if h, err := imageSpec.Digest(); err != nil {
		log.Errorf("get remote image %s digest error %s", ref.Name(), err)
		return "", err
	} else {
		return h.String(), nil
	}
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

		latestTag, err := c.getLatestImageTag(tag)
		if err != nil {
			return err
		}

		inspect, _, err := dockerClient.ImageInspectWithRaw(cmd.Context(), container.ImageID)
		if err != nil {
			return err
		}

		if !c.equalDigest(latestTag.digest, inspect) {
			cmd.Printf("Container %s newer tag \"%s:%s\"\n", container.Image, latestTag.name, latestTag.digest)
		} else {
			log.Debugf("container %s using latest image", container.Image)
		}
	}

	return nil
}

func (c *peekLatestImageCmd) equalDigest(digest string, inspect dockertypes.ImageInspect) bool {
	if len(inspect.RepoDigests) == 0 {
		return false
	}

	for _, repoDigest := range inspect.RepoDigests {
		parts := strings.SplitN(repoDigest, "@", 2)
		if len(parts) == 2 && parts[1] == digest {
			return true
		}
	}

	return false
}

func init() {
	cmd := peekLatestImageCmd{
		timeout: time.Minute,
	}
	realCmd := &cobra.Command{
		Use:   "peek_latest_image",
		Short: "check all containers' latest image",
		RunE:  cmd.RunE,
	}
	rootCmd.AddCommand(realCmd)
	realCmd.Flags().DurationVar(&cmd.timeout, "timeout", cmd.timeout, "timeout")
	realCmd.Flags().VarP(&cmd.tagRegexp, "tag-filter", "f", "regexp for registry tag filtering")
}
