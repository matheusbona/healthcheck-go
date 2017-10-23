package main

import (
   "fmt"
   "net/http"
   "log"
   "strconv"
   "time"
   "encoding/json"
   "os"
   "github.com/sendgrid/sendgrid-go"
   "github.com/sendgrid/sendgrid-go/helpers/mail"
   "github.com/gorilla/mux"
)

// Variáveis configuráveis
var url string = "http://www.ecloudc.com.br"
var threshould int = 5
var intervaloConsulta time.Duration = 60 * time.Second
var timeout time.Duration = 5 * time.Second
// Variáveis configuráveis - Edite até aqui :)

var out string
var fail int = 0
var success int = 0
var dispararNotificacao int = 0
var logs []string
var i int = 0

type error interface {
    Error() string
}

func validaEndpoint() string {
  clienthttp := http.Client { Timeout: timeout }

  in := []byte(`{}`)
  var raw map[string]interface{}
  json.Unmarshal(in, &raw)

  time_start := time.Now()
  response, err := clienthttp.Get(url)
  time_end := time.Now()
  duracao := time_end.Sub(time_start)

  if err != nil {
    fail++
    success = 0
    raw["url"] = url
    raw["status"] = err.Error()
    raw["tempo"] = duracao.Seconds()
    raw["hora"] = time_start.String()
    fmt.Println(err.Error())
    out, _ := json.Marshal(raw)
    fmt.Println(out)
    return string(out)
  } else {
    if (response.StatusCode == 200) {
      success++
      fail = 0
      defer response.Body.Close()
      raw["url"] = url
      raw["status"] = strconv.Itoa(response.StatusCode)
      raw["tempo"] = duracao.Seconds()
      raw["data"] = time_start.String()
      out, _ := json.Marshal(raw)
      return string(out)
    } else { 
      fail++
      success = 0
      raw["url"] = url
      raw["status"] = err.Error()
      raw["tempo"] = duracao.Seconds()
      raw["hora"] = time_start.String()
      fmt.Println(err.Error())
      out, _ := json.Marshal(raw)
      fmt.Println(out)
      return string(out)
    } 
  }

  return string(out)
}

func disparaEmail(nome string, email string, resposta string) {
  from := mail.NewEmail("Monitoramento", "no-reply@monitor.com.br")
  subject := "[ALERTA] Indisponibilidade de Ambiente - Down"
  to := mail.NewEmail(nome, email)
  plainTextContent := "Atenção! Ambiente Probe Status: Down (Tentativas: " + strconv.Itoa(threshould) + ") - Resposta: " + resposta + " - URL: " + url
  htmlContent := "<h1>Atenção!</h1> <h2>Alerta de Ambiente - Probe Status: Down</h2> <br /> <h3>Resposta: " + resposta + "</h3> <h3>URL: " + url + "</h3> <h3> Tentativas: " + strconv.Itoa(threshould) + "</h3>"
  message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
  client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
  response, err := client.Send(message)
  
  if err != nil {
    log.Println(err)
  } else {
    fmt.Println("[Status SendGrid - " + strconv.Itoa(response.StatusCode) + "] - Alerta disparado para " + nome + " - " + email)
  }
}

func ControleProbe() {
  if (os.Getenv("URL_MONITORAMENTO") != ``) {
    url = os.Getenv("URL_MONITORAMENTO")
  } else {
    if (len(os.Args) > 1) {
      url = os.Args[1]
    }
  }

  if (os.Getenv("TIMEOUT_MONITORAMENTO") != ``) {
    valortimeout,err := strconv.Atoi(os.Getenv("TIMEOUT_MONITORAMENTO"))
    if (err != nil) {
      log.Fatal("Falha na conversão do parâmetro TIMEOUT_MONITORAMENTO")
    } else {
      timeout = time.Duration(time.Duration(valortimeout) * time.Second)  
    }
  }

  if (os.Getenv("INTERVALO_MONITORAMENTO") != ``) {
    valorintervalo,err := strconv.Atoi(os.Getenv("INTERVALO_MONITORAMENTO"))
    if (err != nil) {
      log.Fatal("Falha na conversão do parâmetro INTERVALO_MONITORAMENTO")
    } else {
      intervaloConsulta = time.Duration(time.Duration(valorintervalo) * time.Second)
    }
  }
  
  fmt.Println("URL configurada: " + url)
  fmt.Println("Timeout HTTP configurado: " + timeout.String())
  fmt.Println("Intervalo de consulta configurado: " + intervaloConsulta.String())

	for {
    resposta := validaEndpoint()
    logs = append(logs,resposta)
    fmt.Println(resposta)

    if(success >= threshould) {
      fail = 0
      dispararNotificacao = 0
      fmt.Println("Probe - UP [" + strconv.Itoa(success) + "]")
    } else { 
      fmt.Println("Probe - Counting - Sucesso: " + strconv.Itoa(success) + " - Fail: " + strconv.Itoa(fail))
    }

    if (fail >= threshould) {
      success = 0
      fmt.Println("Probe - Down [" + strconv.Itoa(fail) + "]")
      if (dispararNotificacao == 0) {
        disparaEmail(`Matheus Bona`, `mateus.bona@gmail.com`, resposta)
        dispararNotificacao = 1
      }
      fail = 0
    }

    i++
    time.Sleep(intervaloConsulta)
  }
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Path: /")
  json.NewEncoder(w).Encode(logs)
}

func main()  {
  fmt.Println("Pingdom-like by Matheus Bona :D - Iniciando...")
  
  route := mux.NewRouter()
  route.HandleFunc("/", HomeHandler)

  go ControleProbe()

  fmt.Println("API escutando na porta 8000")
  log.Fatal(http.ListenAndServe(":8000", route))
}
