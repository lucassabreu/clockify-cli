package api

import (
	"sync"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// ConvertIntoFullTimeEntry converts a dto.TimeEntryImpl into a dto.TimeEntry
func (c *Client) ConvertIntoFullTimeEntry(tei dto.TimeEntryImpl) (dto.TimeEntry, error) {
	wg := sync.WaitGroup{}
	errChan := make(chan error)
	var err error

	wg.Add(4)

	var p *dto.Project
	go func() {
		if tei.ProjectID == "" {
			wg.Done()
			return
		}

		p, err = c.GetProject(GetProjectParam{
			Workspace: tei.WorkspaceID,
			ProjectID: tei.ProjectID,
		})

		if err != nil {
			errChan <- err
		}

		wg.Done()
	}()

	var u *dto.User

	go func() {
		u, err = c.GetUser(tei.UserID)

		if err != nil {
			errChan <- err
		}

		wg.Done()
	}()

	var task *dto.Task

	go func() {
		if tei.TaskID == "" {
			wg.Done()
			return
		}

		task, err = c.GetTask(GetTaskParam{
			Workspace: tei.WorkspaceID,
			TaskID:    tei.TaskID,
		})

		if err != nil {
			errChan <- err
		}

		wg.Done()
	}()

	tags := make([]dto.Tag, len(tei.TagIDs))

	go func() {
		twg := sync.WaitGroup{}
		twg.Add(len(tei.TagIDs))

		tagsChan := make(chan dto.Tag, len(tei.TagIDs))

		for _, tID := range tei.TagIDs {
			go func(id string) {
				t, err := c.GetTag(GetTagParam{
					Workspace: tei.WorkspaceID,
					TagID:     id,
				})

				if err != nil {
					errChan <- err
				}

				if t != nil {
					tagsChan <- *t
				}

				twg.Done()
			}(tID)
		}

		twg.Wait()
		close(tagsChan)

		i := 0
		for t := range tagsChan {
			tags[i] = t
		}

		wg.Done()
	}()

	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case err = <-errChan:
		return dto.TimeEntry{}, err
	case <-done:
		t := dto.TimeEntry{
			ID:            tei.ID,
			ProjectID:     tei.ProjectID,
			Billable:      tei.Billable,
			Description:   tei.Description,
			IsLocked:      tei.IsLocked,
			TimeInterval:  tei.TimeInterval,
			WorkspaceID:   tei.WorkspaceID,
			TotalBillable: 0,

			Tags: tags,
		}

		if p != nil {
			t.Project = p
			t.HourlyRate = p.HourlyRate
		}

		if u != nil {
			t.User = u
		}

		if task != nil {
			t.Task = task
		}

		return t, nil
	}
}
