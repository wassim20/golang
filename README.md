+++++++++ use this to create a swagger file +++++++++
.\swag.exe init --parseInternal --parseDependency --parseDepth 1

+++++++++ get the library to compile automatically +++++++++
go get github.com/githubnemo/CompileDaemon

+++++++++ use this to compile project automatically +++++++++
CompileDaemon -command="./labs"
