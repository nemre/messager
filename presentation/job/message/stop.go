package message

func (j *job) Stop() {
	if !j.start {
		return
	}

	close(j.stop)
	j.start = false
	j.wg.Wait()
}
