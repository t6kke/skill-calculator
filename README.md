![Genral Tests](https://github.com/t6kke/skill-calculator/actions/workflows/ci_general_tests.yaml/badge.svg) ![Style Tests](https://github.com/t6kke/skill-calculator/actions/workflows/ci_style_tests.yaml/badge.svg) ![Security Tests](https://github.com/t6kke/skill-calculator/actions/workflows/ci_gosec.yaml/badge.svg)

# Skill Calculator

Web solution for serving my [Badminton Skill Calculator](https://github.com/t6kke/BadmintonSkillCalculator) CLI solution.

# Tools Used

- go version 1.24.1
- Styling from [HTML5 Templates](https://html5-templates.com/)
- staticcheck for go linting
- gosec for basic go security tests

# How to Use

**TBD**

end goal is to have this as standalone Docker image that you can just use or you can pull this code build the image yourself

initial docker startup command when image is created,
```
docker run -d -e PLATFORM="<environment>" \
-e PORT="8080" \
-e WEB_ROOT="/var/www/sc/web_assets" \
-e DB_DIR="<example_dir>" \
-e JWT_SECRET="<local_jwt_secret>" \
-p 8080:8080 <image_name:tag> 
```

# Change log

## v7

- New internal go package 'bsc' created to handle interactions with
- Basic unit tests for bsc pacakge created
- Initial bsc package function created for executing BSC with custom parameters
- Excel file upload functionality created; goserver end point, app.js funtion for using it and BSC output response displayed on the website
- Getting league standings functionality created; goserver end point, app.js function for using it and builds tables on website to display rankings
- Getting tournaments list functionality created; goserver end point, app.js functionality to use it and disply it as a list
- Getting tournament result functionality created; goserver end point, app.js functionality to use it and display the results
- Listing categegories functionality created; goserver end point, app.js functionality to use it and list the data
- Adding categorory functionality created; goserver end point, app.js and html form to use it
- Modified how db environment variable works, it now defines directory for databases and webserver has it's own hardcoded db name
- Added Dockerfile for creating local docker image of application, including BSC cli tool
- website JWT expiry validation before useage

## v6

- Unique database name generator created into internal database pacakge
- Internal database package updated with league creation functionality that creates entry into leagues table and to the users_leagues relation table
- Handler function create to enable leagues creation, go server endpoint created to call this function
- Website app.js function created to use leagues creation endpoint
- Internal database package updated with retrieval of specific league based on league id
- Handler function created to use the new database package function and go server endpoint created for it
- Website app.js function created to use league retrieval endpoint
- Internal database package updated with all leagues retrieval for specific user
- Handler function created to use it and go server endpoint created for it
- Website app.js function created to us the all leagues retrieval
- Internal database package updated for league deletion and it also removes the user league relation
- Handler function created to use leage deletion and go server endpoint created to use it
- Website app.js function created to use leage deletion endpoint
- Web UI elements created for leagus management
- Main fuctionality prepared of leages managment
- Removed the content.html logic and main page content is on the index.html page, fixed session and token handling issue
- Fixed issue with login response showing hashed password
- Added logout functionality

## v5

- Fixed app.js "getElementById" errors
- Added basic forms formatting css logic
- Added get token function to go auth package and unit tests for it
- Initial prep for leage creation functionality

## v4

- No password is given back as ouput for user on signup
- Internal auth package created for various authorization workflows
- Password hashing added in auth package, unit tests for it created
- JWT creation and validation added to auth package, unit tests created for it
- Login end point create for server
- Login functionality created in app.js and enabled for login button and on sign up fuction
- Unit tests added to ci workflows
- Database leagues table creation fixed for id field autoincrement

## v3

- Internal database package created for various database interactions
- Initializing and creating new SQLite database file if needed
- Create user to db and get user by e-mail functionality
- New user creation endpoint created on go server
- app.js for frontend created to interact with new endpoint
- Initial DB setup done, users can sign up on website

## v2

- Initial ci automated testing for github workflows
- Resolved some errors brought up by the testing results

## v1

- Initial fileserver for website
- Basic website structure
