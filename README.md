# goclm
A simple incomplete Golang sdk for BMC Cloud Lifecycle Management Product - CLM

Functions implemented:
func Auth(url, userw, passw string) (ClmService, error) - 
    Authenticate a CLM user/password against a CLM URL that is an API end point 
func (c *ClmService) ServiceCreate(offeringName, offeringREID, offeringGUID, serviceName string, quantity int, userName, password, hostNamePrefix, tenantName string) (error) -
    Request an offering
	To Do: Convert ServiceOffering to REID and GUID automatically in the code and remove these 2 args
	
To Do:
func (c *ClmService) ServiceDecommission()
func (c *ClmService) ServiceStart()
func (c *ClmService) ServiceStop()

type ClmService struct {
	authtoken	string
	User		string
	URL			string
}

Sample code:
	...
	clmservice, err := Auth(clmURL, clmUserName, clmUserPassword)
	if err != nil {
		return nil, err
	}
	...
	err:=clmService.ServiceCreate(serviceOffering, serviceOfferingREID, serviceOfferingGUID, serviceName, quantity, userName, userPassword, hostNamePrefix, tenantName)

