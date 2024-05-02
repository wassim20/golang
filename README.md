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
3. Database will migrate autoomatically just configure the port and password in (`.env`) file and create the database with the name of ur choosing.
4. Install packages ex Git :```go mod download``` build and run the backend ex Git :```CompileDaemon -command="./labs"```
5. Install Angular CLI.
6. Navigate to the frontend directory and run ex Git :```npm install``` followed by ex Git :```ng serve -o```.


## Usage

1. Access the application at [http://localhost:4200].
2. Log in with root ex Git :```go run main.go -root``` or create an account to start using the application.
3. Create campaigns, set up automations, and track analytics.

## License

This project is licensed under the [MIT License](LICENSE).

## Contact

For any inquiries, please contact [email@example.com](mailto:email@example.com).

---

Feel free to customize the README.md to fit your project's specific details and requirements.
