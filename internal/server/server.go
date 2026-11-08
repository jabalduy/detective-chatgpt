package server

import (
	"net/http"
	"os"
)

// ЗАПУСК СЕРВЕРА
func StartServer() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/train", trainHandler)
 http.HandleFunc("/examinations", examinationsHandler)
 http.HandleFunc("/examinations/start", examinationStartHandler)
 http.HandleFunc("/evidence", evidenceHandler)
 http.HandleFunc("/evidence/present", evidencePresentHandler)
	http.HandleFunc("/scene", sceneHandler)
	http.HandleFunc("/conversation/start", conversationStartHandler)
	http.HandleFunc("/conversation", conversationHandler)
	http.HandleFunc("/conversation/choose", conversationChooseHandler)
	http.HandleFunc("/characters", charactersHandler)
	http.HandleFunc("/scene/choose", sceneChooseHandler)
	http.HandleFunc("/dialogue/finish", finishDialogueHandler)
	http.HandleFunc("/locations", locationsHandler)
	http.HandleFunc("/objects", objectsHandler)
	http.HandleFunc("/inspect-object", inspectObjectHandler)
	http.HandleFunc("/suspects", suspectsHandler)
	http.HandleFunc("/dialogues", dialoguesHandler)
	http.HandleFunc("/ask-dialogue", askDialogueHandler)
	http.HandleFunc("/case", caseHandler)
	http.HandleFunc("/accuse", accuseHandler)
	http.HandleFunc("/cases", casesHandler)
	http.HandleFunc("/deduce", deductionHandler)
	http.HandleFunc("/facts", factsHandler)
	http.HandleFunc("/motives", motivesHandler)
	http.HandleFunc("/dialogue/start", startDialogueHandler)
	http.HandleFunc("/confront", confrontationHandler)
	http.HandleFunc("/event", eventHandler)
	http.HandleFunc("/event/choose", chooseOptionHandler)
	http.HandleFunc("/event/expire", expireEventHandler)

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
