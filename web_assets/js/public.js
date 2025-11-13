document.addEventListener('DOMContentLoaded', async () => {
    await getLeagues();
});

const leagueStateHandler = createLeagueStateHandler();

function createLeagueStateHandler() {
    let currentLeagueID = null;

    return async function handleLeagueClick(leagueID) {
        if (currentLeagueID !== leagueID) {
            currentLeagueID = leagueID;
            document.getElementById('bsc-response').innerHTML = '';
            document.getElementById('tournament-result-report').innerHTML = '';
            await getLeague(leagueID);
            await getLeageStandings(leagueID);
            await getTournaments(leagueID);
        }
    };
}

async function getLeagues() {
    try {
        const res = await fetch('/api/public_leagues', {
            method: 'GET',
        });
        if (!res.ok) {
            const data = await res.json();
            throw new Error(`Failed to get leagues. Error: ${data.error}`);
        }

        const leagues = await res.json();
        const leagueList = document.getElementById('league-list');
        leagueList.innerHTML = '';
        for (const league of leagues) {
            const listItem = document.createElement('li');
            listItem.textContent = league.title;
            listItem.onclick = () => leagueStateHandler(league.id);
            leagueList.appendChild(listItem);
        }
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}


async function getLeague(leagueID) {
    try {
        const res = await fetch(`/api/leagues/${leagueID}`, {
            method: 'GET',
            headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
        });
        if (!res.ok) {
            throw new Error('Failed to get league.');
        }

        const league = await res.json();
        viewLeague(league);
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}

async function getLeageStandings(leagueID) {
    //TODO standings should be hidden if no results exist
    try {
        const res = await fetch(`/api/public_leagues_standings/${leagueID}`, {
            method: 'GET',
        });
        if (!res.ok) {
            throw new Error('Failed to get Leage standings information.');
        }

        //TODO cleare the table before before adding data in
        const results = await res.json();
        const resultsList = document.getElementById('standings');
        resultsList.innerHTML = '';
        for (const category of results.category) {
            const listItem = document.createElement('p');
            listItem.textContent = category.name + " - " + category.description;
            resultsList.appendChild(listItem);

            const resultTable = document.createElement('table');
            resultTable.setAttribute("class", "table table-bordered table-striped")
            resultsList.appendChild(resultTable);
            const th = document.createElement('thead')
            const tr = document.createElement('tr');
            const tableHeadName = document.createElement('th');
            const tableHeadElo = document.createElement('th');
            tableHeadName.textContent = "Name";
            tableHeadElo.textContent = "ELO";
            resultTable.appendChild(th);
            resultTable.appendChild(tr);
            resultTable.appendChild(tableHeadName);
            resultTable.appendChild(tableHeadElo);
            const tb = document.createElement('tbody');
            resultTable.appendChild(tb);

            for (const rank of category.ranking) {
                const tableRow = document.createElement('tr');
                const tableItemName = document.createElement('td');
                const tableItemElo = document.createElement('td');
                tableItemName.textContent = rank.name;
                tableItemElo.textContent = rank.elo;
                resultTable.appendChild(tableRow);
                resultTable.appendChild(tableItemName);
                resultTable.appendChild(tableItemElo);
            }
        }
    } catch (error) {
        //alert(`Error: ${error.message}`);
        console.log(`Error: ${error.message}`)
    }
}

let currentLeague = null;

function viewLeague(league) {
    currentLeague = league;
    document.getElementById('league-display').style.display = 'block';
    document.getElementById('league-title-display').textContent = league.title;
    document.getElementById('league-description-display').textContent = league.description;
}

async function getTournaments(leagueID) {
    try {
        const res = await fetch(`/api/public_tournamnets/${leagueID}`, {
            method: 'GET',
        });
        if (!res.ok) {
            const data = await res.json();
            throw new Error(`Failed to get leagues. Error: ${data.error}`);
        }

        const bsc_response = await res.json();
        const tournamentsList = document.getElementById('tournamnet-list');
        tournamentsList.innerHTML = '';
        for (const tournament of bsc_response.tournaments) {
            const listItem = document.createElement('li');
            listItem.textContent = tournament.name + " --- date: " + tournament.date + " --- location: " + tournament.location;
            listItem.onclick = () => getTournamentResults(leagueID, tournament.id);
            tournamentsList.appendChild(listItem);
        }
    } catch (error) {
        //alert(`Error: ${error.message}`);
        console.log(`Error: ${error.message}`)
    }
}

async function getTournamentResults(leagueID, tournamentID) {
    try {
        const res = await fetch(`/api/public_tournamnets/${leagueID}/${tournamentID}`, {
            method: 'GET',
        });
        if (!res.ok) {
            const data = await res.json();
            throw new Error(`Failed to get leagues. Error: ${data.error}`);
        }

        const bsc_response = await res.json();
        const resultsList = document.getElementById('tournament-result-report');
        resultsList.innerHTML = '';
        const resultTable = document.createElement('table');
        resultTable.setAttribute("class", "table table-bordered table-striped")
        resultsList.appendChild(resultTable);
        const th = document.createElement('thead');
        const tr = document.createElement('tr');
        const tableHeadPosition = document.createElement('th');
        const tableHeadTeam = document.createElement('th');
        const tableHeadGamesPlayed = document.createElement('th');
        const tableHeadGamesWon = document.createElement('th');
        const tableHeadPointsFor = document.createElement('th');
        const tableHeadPointsAgainst = document.createElement('th');
        const tableHeadPointsDiff = document.createElement('th');
        tableHeadPosition.textContent = "Position";
        tableHeadTeam.textContent = "Team";
        tableHeadGamesPlayed.textContent = "Games Played";
        tableHeadGamesWon.textContent = "Games Won";
        tableHeadPointsFor.textContent = "Points For";
        tableHeadPointsAgainst.textContent = "Points Against";
        tableHeadPointsDiff.textContent = "Points Difference";
        resultTable.appendChild(th);
        resultTable.appendChild(tr);
        resultTable.appendChild(tableHeadPosition);
        resultTable.appendChild(tableHeadTeam);
        resultTable.appendChild(tableHeadGamesPlayed);
        resultTable.appendChild(tableHeadGamesWon);
        resultTable.appendChild(tableHeadPointsFor);
        resultTable.appendChild(tableHeadPointsAgainst);
        resultTable.appendChild(tableHeadPointsDiff);
        const tb = document.createElement('tbody');
        resultTable.appendChild(tb);
        for (const result of bsc_response.results) {
            const tableRow = document.createElement('tr');
            const tableItemPosition = document.createElement('td');
            const tableItemTeam = document.createElement('td');
            const tableItemGamesPlayed = document.createElement('td');
            const tableItemGamesWon = document.createElement('td');
            const tableItemPointsFor = document.createElement('td');
            const tableItemPointsAgainst = document.createElement('td');
            const tableItemPointsDiff = document.createElement('td');
            tableItemPosition.textContent = result.position;
            tableItemTeam.textContent = result.team;
            tableItemGamesPlayed.textContent = result.games_total;
            tableItemGamesWon.textContent = result.games_won;
            tableItemPointsFor.textContent = result.points_for;
            tableItemPointsAgainst.textContent = result.points_against;
            tableItemPointsDiff.textContent = result.points_diff;
            resultTable.appendChild(tableRow);
            resultTable.appendChild(tableItemPosition);
            resultTable.appendChild(tableItemTeam);
            resultTable.appendChild(tableItemGamesPlayed);
            resultTable.appendChild(tableItemGamesWon);
            resultTable.appendChild(tableItemPointsFor);
            resultTable.appendChild(tableItemPointsAgainst);
            resultTable.appendChild(tableItemPointsDiff);
        }
        document.getElementById('tournament-result-container').style.display = 'block';
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}
