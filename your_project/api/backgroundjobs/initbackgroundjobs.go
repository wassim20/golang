package backgroundjobs

func InitBackgroundJob() {
	err := StartBackgroundJob("@every 10s") // Adjust schedule for testing
	if err != nil {
		panic(err) // Or handle errors more gracefully
	}
}
