#base image from DockerHub, python base image used because of BSC
FROM python:3.12

# install additional depndent packages for BSC
RUN pip install pandas
RUN pip install openpyxl

# install git for pulling BSC
RUN apt-get update && \
    apt-get install -y git

# clone BSC to image
#RUN git clone https://github.com/t6kke/BadmintonSkillCalculator.git /opt/BSC
RUN git clone --depth 1 --branch alpha4 https://github.com/t6kke/BadmintonSkillCalculator.git /opt/BSC

# remove git
RUN apt-get remove -y git

# copy over go server
COPY skill-calculator /bin/skill-calculator

# copy over web assets
COPY ./web_assets /var/www/sc/web_assets

# runs the server software
CMD ["/bin/skill-calculator"]
