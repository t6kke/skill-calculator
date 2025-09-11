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

# Change log

## v6

- Unuque database name generator created into internal database pacakge
- Internal database package updated with League creattion, reteival by id and retreival for all based on currentl user
- API endpoinds created for Leagues, for creation, retreival with specific id, retreival for specific user for all leagues
- Website app.js fucntion to use new endpoint for creation of the leage
- Main fuctionality of leages managment completed
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
