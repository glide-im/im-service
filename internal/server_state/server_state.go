package server_state

import (
	"encoding/json"
	"fmt"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/im-service/internal/im_server"
	"io"
	"net/http"
)

type httpHandler struct {
	g *im_server.GatewayServer
}

func (h *httpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			bytes := []byte(fmt.Sprintf("%v", err))
			logger.E("%v", err)
			_, _ = writer.Write(bytes)
		}
	}()
	if request.URL.Path == "/state" {
		writer.WriteHeader(200)
		state := h.g.GetState()
		bytes, err2 := json.Marshal(&state)
		if err2 != nil {
			panic(err2)
		}
		_, err := writer.Write(bytes)
		if err != nil {
			panic(err)
		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func StartSrv(port int, server *im_server.GatewayServer) {
	a := fmt.Sprintf("0.0.0.0:%d", port)
	err := http.ListenAndServe(a, &httpHandler{g: server})
	if err != nil {
		return
	}
}

func ShowServerState(addr string) {
	url := fmt.Sprintf("http://%s/state", addr)
	fmt.Printf("get state from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("get state from %s failed, status code %d\n", url, resp.StatusCode)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if len(bytes) <= 0 {
		fmt.Println("get state failed, empty body")
		return
	}

	var state im_server.GatewayMetrics
	err = json.Unmarshal(bytes, &state)
	if err != nil {
		panic(err)
	}
	printBeautifulState(&state)
}

func printBeautifulState(state *im_server.GatewayMetrics) {
	fmt.Printf("\nServerId:\t%s", state.ServerId)
	fmt.Printf("\nAddr:\t\t%s", state.Addr)
	fmt.Printf("\nPort:\t\t%d", state.Port)
	fmt.Printf("\nStartAt:\t%s", state.StartAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("\nRunningHours:\t\t%.2f", state.RunningHours)

	fmt.Println("== metric ==")
	fmt.Printf("\nOnlineClients:\t\t%d", state.Conn.ConnectionCounter.Count())
	fmt.Printf("\nOnlineTempClients:\t%d", state.Conn.OnlineTempCounter.Count())
	fmt.Printf("\nTempConnAliveMaxSec:\t\t%d", state.Conn.AliveTempH.Max())
	fmt.Printf("\nTempConnAliveMeanSec:\t\t%f", state.Conn.AliveTempH.Mean())
	fmt.Printf("\nOnlineTempClients:\t%d", state.Conn.OnlineTempCounter.Count())
	fmt.Printf("\nLoggedConnAliveMaxSec:\t\t%d", state.Conn.AliveLoggedH.Max())
	fmt.Printf("\nLoggedConnAliveMeanSec:\t\t%f", state.Conn.AliveLoggedH.Mean())
	fmt.Printf("\nOutMessages:\t%d", state.Message.OutCounter.Count())
	fmt.Printf("\nOutMessageFails:\t%d", state.Message.FailsCounter.Count())
	fmt.Printf("\nInMessages:\t%d", state.Message.InCounter.Count())
	fmt.Printf("\n")
}
