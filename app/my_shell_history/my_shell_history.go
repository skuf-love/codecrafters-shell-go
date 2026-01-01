package my_shell_history

var log []string

func StoreCommand(cmd string){
	log = append(log, cmd)
}

func Log() []string {
	return log
}
