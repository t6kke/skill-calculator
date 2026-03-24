#base image from DockerHub, python base image used because of BSC
FROM python:3.12

# install libraries required by Chrome
RUN apt-get -y update && apt-get install -y libnspr4 libnss3-tools

# setting up chrome for selenium
# Adding trusting keys to apt for repositories
RUN wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | tee /etc/apt/trusted.gpg.d/google.asc >/dev/null
# Adding Google Chrome to the repositories
RUN sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list'
# Updating apt and install Google Chrome
RUN apt-get -y update && apt-get install -y google-chrome-stable

# install additional depndent packages for BSC
RUN pip install pandas
RUN pip install openpyxl
RUN pip install beautifulsoup4
RUN pip install selenium
RUN pip install --upgrade requests

# install git for pulling BSC
RUN apt-get update && \
    apt-get install -y git

# clone BSC to image
#RUN git clone https://github.com/t6kke/BadmintonSkillCalculator.git /opt/BSC
RUN git clone --depth 1 --branch alpha5 https://github.com/t6kke/BadmintonSkillCalculator.git /opt/BSC

# remove git
RUN apt-get remove -y git

# copy over go server
COPY skill-calculator /bin/skill-calculator

# copy over web assets
COPY ./web_assets /var/www/sc/web_assets

# runs the server software
CMD ["/bin/skill-calculator"]
