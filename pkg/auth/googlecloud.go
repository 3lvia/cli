package auth

func AuthenticateGoogle() error {
	/*
		gcloudAuthLoginOutput := gcloudAuthLoginCommand(nil)
		if command.IsError(gcloudAuthLoginOutput) {
			return fmt.Errorf("Failed to authenticate to Google Cloud: %w", gcloudAuthLoginOutput.Error)
		}
	*/

	return nil
}

/*
func gcloudAuthLoginCommand(
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"gcloud",
			"auth",
			"login",
		),
		runOptions,
	)
}
*/
