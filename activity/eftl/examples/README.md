
##Setup
Step 1 : Install FTL Server
a. FTL download(/https://www.tibco.com/products/tibco-ftl);
b. Follow the installation instructions for your platform(https://docs.tibco.com/pub/ftl/5.3.2/doc/pdf/TIB_ftl_5.3_Installation.pdf)


Step 2: Install eFTL Server
a. EFTL download(https://www.tibco.com/products/tibco-eftl);
b. Follow the installation instructions for your platform here(https://docs.tibco.com/pub/eftl/3.2.0/doc/html/GUID-9F5E7521-39B1-4DFD-B2E6-35164F9406CD.html)


Step 3: To start the EFTL server run: go run helper/main.go -ftl

##Then in another terminal run: go run helper/main.go -eftl


Step 4: To install run the following commands:

flogo create -f flogo.json
cd eftl
flogo build

#Run
bin/eftl
This sends message to eftl server

#Testing
Get client file
"git clone github.com/project-flogo/contrib/activity/eftl"

Now run the client file:
go run client/client.go

Open a new terminal to call activity:
curl -d "{\"message\": \"hello world\"}" -X GET http://localhost:9096/eftl

Client serves a demo eftl server.You should then see something like on client screen:
{"message": "hello world"}


