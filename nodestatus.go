package main
import (
    "bufio"
    "encoding/json"
    "fmt"
    "net/http"
    "os/exec"
    "regexp"
    "strings"
    "time"
)

func openPorts() map[string]string {
    data := make(map[string]string)
    output, err := exec.Command("netstat", "-ltn").Output()
    if err != nil {
        // fmt.Println(err)
    } else {
        portRegex := regexp.MustCompile(`:([0-9]+)\s`)
        scanner := bufio.NewScanner(strings.NewReader(string(output)))
        for scanner.Scan() {
            line := scanner.Text()
            if (strings.Contains(line, ":") ) {
               // fmt.Print(line + " .... ")
               port := strings.TrimSpace(strings.Replace(portRegex.FindString(line), ":", "", -1))
               // fmt.Println(port)
               data[port] = port
            }
        }
    }
    return data
}

var ports map[string]string = openPorts()
func portStatus(p int) string {
    _, present := ports[fmt.Sprintf("%d", p)]
    return fmt.Sprint(present)
}

func urlStatus(url string) string {
    timeout := time.Duration(3 * time.Second)
    client := http.Client{ Timeout: timeout }
    response, err := client.Get(url)
    if err != nil {
        return fmt.Sprint(err)
    } else {
        defer response.Body.Close()
        return fmt.Sprintf("%d", response.StatusCode)
    }
}

type StatusItem struct {
    Label string  `json:"label"`
    Status string `json:"status"`
}

func getProps() map[string]string  {
    return map[string]string {
      "name": "localhost",
    } 
}

func main() {
    props := getProps()
    statuses := []StatusItem{
        StatusItem{ "Secure Port (SSL)", portStatus(443)}, 
        StatusItem{ "MySQL Service", portStatus(3306) },
        StatusItem{ "Local CouchDB", portStatus(5984)}, 
        StatusItem{ "Redis KV Store", portStatus(6379) }, 
        StatusItem{ "Front HTTP (" + props["name"] + ")", urlStatus("http://" + props["name"] + "/") },
        StatusItem{ "Front HTTP (" + props["name"] + ")", urlStatus("https://" + props["name"] + "/") },
    }
    json, _ := json.Marshal(statuses)
    fmt.Println(string(json))
}
