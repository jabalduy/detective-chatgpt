package server

import (
	"net/http"
	"os"
)

// ЗАПУСК СЕРВЕРА
func NewHandler() http.Handler {
 mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/start", startHandler)
 mux.HandleFunc("/save", saveHandler)
 mux.HandleFunc("/load", loadHandler)
	mux.HandleFunc("/time", timeHandler)
	mux.HandleFunc("/train", trainHandler)
 mux.HandleFunc("/examinations", examinationsHandler)
 mux.HandleFunc("/examinations/start", examinationStartHandler)
 mux.HandleFunc("/evidence", evidenceHandler)
 mux.HandleFunc("/evidence/present", evidencePresentHandler)
	mux.HandleFunc("/scene", sceneHandler)
	mux.HandleFunc("/conversation/start", conversationStartHandler)
	mux.HandleFunc("/conversation", conversationHandler)
	mux.HandleFunc("/conversation/choose", conversationChooseHandler)
	mux.HandleFunc("/characters", charactersHandler)
	mux.HandleFunc("/scene/choose", sceneChooseHandler)
	mux.HandleFunc("/dialogue/finish", finishDialogueHandler)
	mux.HandleFunc("/locations", locationsHandler)
	mux.HandleFunc("/objects", objectsHandler)
	mux.HandleFunc("/inspect-object", inspectObjectHandler)
	mux.HandleFunc("/suspects", suspectsHandler)
	mux.HandleFunc("/dialogues", dialoguesHandler)
	mux.HandleFunc("/ask-dialogue", askDialogueHandler)
	mux.HandleFunc("/case", caseHandler)
	mux.HandleFunc("/accuse", accuseHandler)
 mux.HandleFunc("/ending", endingHandler)
	mux.HandleFunc("/cases", casesHandler)
	mux.HandleFunc("/deduce", deductionHandler)
	mux.HandleFunc("/facts", factsHandler)
	mux.HandleFunc("/motives", motivesHandler)
	mux.HandleFunc("/dialogue/start", startDialogueHandler)
	mux.HandleFunc("/confront", confrontationHandler)
	mux.HandleFunc("/event", eventHandler)
	mux.HandleFunc("/event/choose", chooseOptionHandler)
	mux.HandleFunc("/event/expire", expireEventHandler)

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

 return saveAfterMutation(mux)
}

func StartServer() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, NewHandler())
	if err != nil {
		panic(err)
	}
}
