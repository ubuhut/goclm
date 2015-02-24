package clm

import (
	"encoding/json"
	"fmt"
	"os"
	"net/http"
	"bytes"
	"io/ioutil"
	"time"
	//log "github.com/Sirupsen/logrus"
)

const (
  AUTHORIZATION_TOKEN	= "Authentication-Token"
  AUTH_ACTION = "/login"
  SVC_SRCH_ACTION = "/ServiceOffering/search"
  SVC_BULK_CREATE_ACTION = "/ServiceOfferingInstance/bulkCreate"

  OP_TYPE_STRING = "java.lang.String"
  OP_TYPE_DATE = "java.util.Date"
  OP_TYPE_INT = "java.lang.Integer"
  OP_TYPE_SVC = "com.bmc.cloud.model.beans.ServiceOffering"
  OP_TYPE_SVC_INST = "com.bmc.cloud.model.beans.ServiceOfferingInstance"
  TIMEOUT = 300
  UP	= 200
  DOWN	= 0
)


type ClmService struct {
	authtoken	string
	User		string
	URL		string
}
type ClmTaskResponse struct {
	Cloudclass string
	ClassName string
	CreationTime string
	Guid string
	OperationName string
	TaskInternalUUID string
	TaskState string
	TaskStatusURI string
	TransactionID string
}


func Auth(url, userw, passw string) (ClmService, error) {
	fmt.Println("Auth:",url)
	c := ClmService{}
	c.URL=url
	url=url+AUTH_ACTION
	fmt.Println("Auth:", url, userw)
	type AuthJson struct {
		Name   string `json:"username"`
		Password	string `json:"password"`
	}
	group := AuthJson{
		Name:  userw ,
		Password: passw,
	}
	b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("Auth: error:", err)
		return  c,err
	}
	os.Stdout.Write(b)
	jsonStr := b
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")
	fmt.Println("Auth: req=%s", req, "\n")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
		fmt.Println("Auth: client.Do req:%s", req, err)
		return c, err
 	}
	var status int
	if (err == nil) && (resp.StatusCode == 200) {
        	status = UP
    } else {
        	status = DOWN
    }
	fmt.Println("Auth: StatusCode: %s status:",resp.StatusCode, status)
    defer resp.Body.Close()
    fmt.Println("Auth: response Status:", resp.Status)
    fmt.Println("Auth: response Headers:", resp.Header)
	contents, err := ioutil.ReadAll(resp.Body)
    if err != nil {
            fmt.Printf("Auth: ioutil.ReadAll err:%s", err)
            return c, err
    }
    fmt.Printf("Auth: response Contents:%s\n", string(contents))

	token := resp.Header[AUTHORIZATION_TOKEN]
	fmt.Println("Auth: token:", token)
	if token == nil {
		fmt.Println("Auth: no auth token receieved")
		c.authtoken = "aaa"
		return c,fmt.Errorf("No auth token")
	}
	c.authtoken = token[0]
	fmt.Println("Auth: Token:", token, c.authtoken)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Auth: Response Body:", string(body))
	return c, nil
}
func (c* ClmService) Authenticate(url, userw, passw string) (error) {
	c.URL=url
	url=url+AUTH_ACTION
	fmt.Println("Auth:", url, userw)
	type AuthJson struct {
		Name   string `json:"username"`
		Password	string `json:"password"`
	}
	group := AuthJson{
		Name:  userw ,
		Password: passw,
	}
	b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("Auth: error:", err)
		return  err
	}
	os.Stdout.Write(b)
	jsonStr := b
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")
	fmt.Println("Auth: req=%s", req, "\n")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
		fmt.Println("Auth: client.Do req:%s", req, err)
		return err
 	}
	var status int
	if (err == nil) && (resp.StatusCode == 200) {
        	status = UP
    } else {
        	status = DOWN
    }
	fmt.Println("Auth: StatusCode: %s status:",resp.StatusCode, status)
    defer resp.Body.Close()
    fmt.Println("Auth: response Status:", resp.Status)
    fmt.Println("Auth: response Headers:", resp.Header)
	contents, err := ioutil.ReadAll(resp.Body)
    if err != nil {
            fmt.Printf("Auth: ioutil.ReadAll err:%s", err)
            return err
    }
    fmt.Printf("Auth: response Contents:%s\n", string(contents))

	token := resp.Header[AUTHORIZATION_TOKEN]
	fmt.Println("Auth: token:", token)
	if token == nil {
		fmt.Println("Auth: no auth token receieved")
		c.authtoken = "aaa"
		return fmt.Errorf("No auth token")

	} else {
		c.authtoken = token[0]
		fmt.Println("Auth: Token:", token, c.authtoken)
    	body, _ := ioutil.ReadAll(resp.Body)
    	fmt.Println("Auth: Response Body:", string(body))
	}
	return nil
}

func (c *ClmService)	GetTask(taskURL string) (ClmTaskResponse, error) {
	fmt.Println("GetTask: taskurl=", taskURL)
	req, err := http.NewRequest("GET", taskURL, nil)
    req.Header.Set(AUTHORIZATION_TOKEN, c.authtoken)
    req.Header.Set("Content-Type", "application/json")
	fmt.Println("GetTask: req=",req)
	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
		fmt.Printf("GetTask: Client.Do err=%s", err)
		return ClmTaskResponse{}, err
	}
    defer resp.Body.Close()
    fmt.Println("GetTask response Status:", resp.Status)
    fmt.Println("GetTask response Headers:", resp.Header)
	contents, err := ioutil.ReadAll(resp.Body)
    if err != nil {
            fmt.Printf("GetTask ioutil %s", err)
			return ClmTaskResponse{}, err
    }
    fmt.Printf("GetTask: response contents:%s\n", string(contents))
	
	var ctaskresp []ClmTaskResponse // TopTracks map[string]interface{}
    err = json.Unmarshal(contents, &ctaskresp)
    if err != nil {
        fmt.Printf("GetTask: json.Unmarshal err:%s", err)
		return ClmTaskResponse{}, err
    }
	//TODO: BUG check size of array
	return ctaskresp[0],nil
}
//Python SDK: obj = service_create(gcac, offeringName=None, serviceName=None, quantity=1, userName=None, password=None, hostNamePrefix=None, tenantName=None, options=None, params=None, decommissiondate=None)

func (c *ClmService) ServiceCreate(offeringName, offeringREID, offeringGUID, serviceName string, quantity int, userName, password, hostNamePrefix, tenantName string) (error) {
	
	fmt.Println("ServiceCreate:", offeringName, offeringREID, serviceName, tenantName, c.authtoken)
	if (c.authtoken == "") {
		fmt.Println("ServiceCreate: error: no auth token!!!, for now still continuing as this is pre-alpha code in test")
		//return nil
	}
	//now setup the json needed to make a REST API call into CLM
	type op struct {
		Name string `json:"name"`
		Value string `json:"value"`
		Optype string `json:"type"`
		Multiplicity string `json:"multiplicity"`
	}
	type op_int struct {
		Name string `json:"name"`
		Value int `json:"value"`
		Optype string `json:"type"`
		Multiplicity string `json:"multiplicity"`
	}
	type t_so struct {
		CloudClass string  `json:"cloudClass"`//: "com.bmc.cloud.model.beans.ServiceOffering",
		Guid string `json:"guid"` //"OI-A6271B0404624169AA0990CBDEF036E5",
		ReconciliationID string `json:"reconciliationID"`  // "76053A5B-70CD-007D-F0AD-3DC42939A8A7"
	}
	type t_op_so2 struct {
		Name string `json:"name"`
		Value t_so `json:"value"`
		Optype string `json:"type"`
		Multiplicity string `json:"multiplicity"`
	}
	op_so := op{
		Name:"serviceOfferingID", 
		Value:offeringREID, 
		Optype:"java.lang.String", 
		Multiplicity:"1",
	}
	/*
	op_sn := op{
		Name:"serviceNameID", 
		Value:serviceName, 
		Optype:"java.lang.String", 
		Multiplicity:"1",
	} */
	sovalue := t_so {
		CloudClass:"com.bmc.cloud.model.beans.ServiceOffering",
		Guid: offeringGUID,
		ReconciliationID: offeringREID,
	}
	op_so2 := t_op_so2{
		Name:"serviceOfferingID", 
		Value:sovalue, 
		Optype:"com.bmc.cloud.model.beans.ServiceOffering", 
		Multiplicity:"1",
	}	
	fmt.Println("ServiceCreate: keep Go Happy",op_so, op_so2)
	op_quantity := op_int{
			Name: "quantity",
			Value: quantity,
			Optype: "java.lang.Integer",
			Multiplicity: "1",
	}
	op_password := op{
			Name: "password",
			Value: password,
			Optype: "java.lang.String",
			Multiplicity: "0..1",
    }
	op_username := op{
			Name: "username",
			Value: userName,
			Optype: "java.lang.String",
			Multiplicity: "0..1",
	}
	op_hostnameprefix := op{
			Name: "hostnamePrefix",
			Value: hostNamePrefix,
			Optype: "java.lang.String",
			Multiplicity: "0..1",
	}
	op_tenant := op{
			Name: "tenant",
			Value: tenantName,
			Optype: "java.lang.String",
			Multiplicity: "0..1",
	}
	op_reqdefinitionname := op{
			Name: "name",
			Value: serviceName,
			Optype: "java.lang.String",
			Multiplicity: "0..1",
	}
	
	fmt.Println("Service-Create: keep go happy name:", op_reqdefinitionname)
	type Any interface{}
	type ServiceRequest2 struct {
		OperationParams [] Any `json:"operationParams"`
		Timeout	int		`json:"timeout"`
		User string `json:"user"`
		PreCallout	string	`json:"preCallout"`
		PostCallout	string	`json:"postCallout"`
        CallbackURL	string	`json:"callbackURL"`
        AlreadyTraversedGlobalRegistry	bool	`json:"alreadyTraversedGlobalRegistry"`
        AlreadyTraversedLocalRegistry bool	`json:"alreadyTraversedLocalRegistry"`
	}
		
	sr2 := ServiceRequest2{}
	sr2.OperationParams = append(sr2.OperationParams, op_so2, op_quantity, op_username, op_password, op_hostnameprefix, op_tenant,op_reqdefinitionname)
	sr2.Timeout = (0)
	sr2.User = userName
	sr2.AlreadyTraversedGlobalRegistry = false
	bsr, err := json.MarshalIndent(sr2,"", "    ")
	if (err != nil) {
		fmt.Println("ServiceCreate: json.MarshalIndent failed with error:%s", err)
		return err
	}
	fmt.Println("ServiceCreate: JSON:")
	os.Stdout.Write(bsr)

	req, err := http.NewRequest("POST", c.URL+"/serviceofferinginstance/bulkcreate", bytes.NewBuffer(bsr))
    req.Header.Set(AUTHORIZATION_TOKEN, c.authtoken)
    req.Header.Set("Content-Type", "application/json")
	fmt.Println("sending...",req)
	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("ServiceCreate: error in client.Do err:%s",err)
		return err
	}
    defer resp.Body.Close()
    fmt.Println("ServiceCreate: response Status:", resp.Status)
    fmt.Println("ServiceCreate: response Headers:", resp.Header)
	contents, err := ioutil.ReadAll(resp.Body)
    if err != nil {
            fmt.Printf("ServiceCreate: error in ioutil.ReadAll body err:%s", err)
            os.Exit(1)
    }
    fmt.Printf("ServiceCreate: response contents:%s\n", string(contents))
	
	type clm_resp struct {
		Cloudclass string
		ClassName string
		CreationTime string
		Guid string
		OperationName string
		TaskInternalUUID string
		TaskState string
		TaskSubState string
		TaskProgress int
		TaskStatusURI string
		TransactionID string
		Errors string
		IsCallout bool
		IsError bool
		IsSuccess bool
	}
	var cresp []clm_resp // TopTracks map[string]interface{}
    err = json.Unmarshal(contents, &cresp)
    if err != nil {
        fmt.Println("ServiceCreate: error in json.Unmarshal err:%s", err)
		return err
    }
    fmt.Printf("ServiceCreate: Results2: %v\n", cresp)
	if (len(cresp) > 0) {
		//TODO: verify if response is good and if we got immediate results or a task URL?
		url := cresp[0].TaskStatusURI
		state :=cresp[0].TaskState
		fmt.Println("ServiceCreate: TaskStatusURI: & State:", url, state)
		var count int
		for (count < 200) {
			count += count
			ctaskresp, err := c.GetTask(url)
			if (err!=nil) {
				fmt.Println("ServiceCreate: error in GetTask err:%s", err)
				return err
			}
			if ((ctaskresp.TaskState == "FAILED")||(ctaskresp.TaskState =="COMPLETED")) {
				fmt.Println("ServiceCreate: FAILED or COMPLETED Task!!! TaskState:%s", ctaskresp.TaskState)
				return nil
			}
			time.Sleep(10 * time.Second)
		}//for
	} else {
		//error in response!
		return fmt.Errorf("No auth token")
	}
	return err	
}

