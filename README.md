# Email Marketing Project

## Description

This project is a mail marketing application that allows users to create campaigns and automations. It is built using Golang 1.19 for the backend and Angular for the frontend. The backend is developed with the Gin framework, and the database used is PostgreSQL.

## Key Features

- **Campaign Creation:** Users can create email campaigns with customized content and settings.
- **Automation:** Users can set up automated email sequences based on triggers and conditions.
- **Analytics:** The application provides analytics for tracking the performance of campaigns.
- **User Management:** Administrators can manage user accounts and permissions.

## Technologies Used

- **Backend:** Golang 1.19, Gin framework
- **Frontend:** Angular
- **Database:** PostgreSQL

## Installation

1. Clone the repository.
2. Install Golang 1.19 and PostgreSQL.
3. Configure the port and password in the `.env` file and create the database with the name of your choosing.
4. Install packages:
    ```bash
    go mod download
    ```
5. Build and run the backend:
    ```bash
    CompileDaemon -command="./labs"
    ```
6. Install Angular CLI.
7. Navigate to the frontend directory and run:
    ```bash
    npm install
    ng serve -o
    ```
8. If you want to use swagger:
    ```bash
    .\swag.exe init --parseInternal --parseDependency --parseDepth 1
    ```
   Then copy the file from `./docs/swagger.json`.

## Usage

1.Log in with root:
    ```bash
    go run main.go -root
    ```
2. Access the application at [http://localhost:4200](http://localhost:4200).

   Or create an account to start using the application.
3. Create campaigns, set up automations, and track analytics.

## License

This project is licensed under the [MIT License](LICENSE).

## Contact

For any inquiries, please contact [email@example.com](mailto:email@example.com).
