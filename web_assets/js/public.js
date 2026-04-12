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

            var resultTable = document.createElement('table');
            resultTable.setAttribute("class", "table table-bordered table-striped")
            resultsList.appendChild(resultTable);
            header = resultTable.createTHead();
            row = resultTable.insertRow(0);
            th0 = document.createElement("th");
            th1 = document.createElement("th");
            th2 = document.createElement("th");
            th0.textContent = "Rank";
            th1.textContent = "Name";
            th2.textContent = "ELO";
            row.appendChild(th0);
            row.appendChild(th1);
            row.appendChild(th2);
            let rank_nbr = 0;
            for (const rank of category.ranking) {
                rank_nbr++;
                row = resultTable.insertRow();
                cell0 = row.insertCell();
                cellA = row.insertCell();
                cellB = row.insertCell();
                cell0.innerHTML = rank_nbr;
                cellA.innerHTML = UpperName(rank.name);
                cellB.innerHTML = rank.elo;
            }
        }
    } catch (error) {
        //alert(`Error: ${error.message}`);
        console.log(`Error: ${error.message}`)
    }
}

function UpperName(name) {
    const exception_words = ["van", "de", "der"]
    let name_parts = name.split(" ");
    for (let i = 0; i < name_parts.length; i++) {
        if (exception_words.includes(name_parts[i])) {
            continue;
        }
        name_parts[i] = name_parts[i].charAt(0).toUpperCase() + name_parts[i].slice(1);
    }
    return name_parts.join(" ");
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

        var resultTable = document.createElement('table');
        resultTable.setAttribute("class", "table table-bordered table-striped")
        resultsList.appendChild(resultTable);
        header = resultTable.createTHead();
        row = resultTable.insertRow(0);
        th1 = document.createElement("th");
        th2 = document.createElement("th");
        th3 = document.createElement("th");
        th4 = document.createElement("th");
        th5 = document.createElement("th");
        th6 = document.createElement("th");
        th7 = document.createElement("th");
        th1.textContent = "Position";
        th2.textContent = "Team";
        th3.textContent = "Games Played";
        th4.textContent = "Games Won";
        th5.textContent = "Points For";
        th6.textContent = "Points Against";
        th7.textContent = "Points Difference";
        row.appendChild(th1);
        row.appendChild(th2);
        row.appendChild(th3);
        row.appendChild(th4);
        row.appendChild(th5);
        row.appendChild(th6);
        row.appendChild(th7);
        for (const result of bsc_response.results) {
            row = resultTable.insertRow();
            cellA = row.insertCell();
            cellB = row.insertCell();
            cellC = row.insertCell();
            cellD = row.insertCell();
            cellE = row.insertCell();
            cellF = row.insertCell();
            cellG = row.insertCell();
            cellA.innerHTML = result.position;
            cellB.innerHTML = result.team;
            cellC.innerHTML = result.games_total;
            cellD.innerHTML = result.games_won;
            cellE.innerHTML = result.points_for;
            cellF.innerHTML = result.points_against;
            cellG.innerHTML = result.points_diff;
        }
        document.getElementById('tournament-result-container').style.display = 'block';
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}
