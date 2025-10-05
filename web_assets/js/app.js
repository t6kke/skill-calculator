document.addEventListener('DOMContentLoaded', async () => {
    const token = localStorage.getItem('token');

    const payload = JSON.parse(atob(token.split('.')[1]));
    const isExpired = payload.exp * 1000 < Date.now();

    if (isExpired) {
        localStorage.removeItem('token'); // clear expired token
        token = null
    }

    //TODO also validate that user from token also exists

    if (token) {
        document.getElementById('auth-section').style.display = 'none';
        document.getElementById('conent-section').style.display = 'block';
        await getLeagues();
    } else {
        document.getElementById('auth-section').style.display = 'block';
        document.getElementById('conent-section').style.display = 'none';
    }
});

document.getElementById('league-draft-form').addEventListener('submit', async (event) => {
    event.preventDefault();
    await createLeague();
});

document.getElementById('login-form').addEventListener('submit', async (event) => {
    event.preventDefault();
    await login();
});

async function signup() {
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    try {
        const res = await fetch('/api/users', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
        });
        if (!res.ok) {
            const data = await res.json();
            throw new Error(`Failed to create user: ${data.error}`);
        }
        console.log('User created!');
        await login();
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}

async function login() {
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    try {
        const res = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
        });
        const data = await res.json();
        if (!res.ok) {
            throw new Error(`Failed to login: ${data.error}`);
        }

        if (data.token) {
            localStorage.setItem('token', data.token);
            document.getElementById('auth-section').style.display = 'none';
            document.getElementById('conent-section').style.display = 'block';
            await getLeagues();
        } else {
            alert('Login failed. Please check your credentials.');
        }
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}

function logout() {
    localStorage.removeItem('token');
    document.getElementById('auth-section').style.display = 'block';
    document.getElementById('conent-section').style.display = 'none';
}

const leagueStateHandler = createLeagueStateHandler();

function createLeagueStateHandler() {
    let currentLeagueID = null;

    return async function handleLeagueClick(leagueID) {
        if (currentLeagueID !== leagueID) {
            currentLeagueID = leagueID;
            await getLeague(leagueID);
            await getLeageStandings(leagueID);
        }
    };
}

async function createLeague() {
    const title = document.getElementById('league-title').value;
    const description = document.getElementById('league-description').value;

    try {
        const res = await fetch('/api/leagues', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify({ title, description }),
        });
        const data = await res.json();
        if (!res.ok) {
            throw new Error(`Failed to create league: ${data.error}`);
        }

        const leagueID = data.id;
        if (leagueID) {
            await getLeagues();
            await leagueStateHandler(leagueID); //TODO do I need this?
        }
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}

async function getLeagues() {
    try {
        const res = await fetch('/api/leagues', {
            method: 'GET',
            headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
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
    try {
        const res = await fetch(`/api/league_standings/${leagueID}`, {
            method: 'GET',
            headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
        });
        if (!res.ok) {
            throw new Error('Failed to get Leage standings information.');
        }

        const league = await res.json();
        //viewLeague(league);
        //TODO logic here to push content to html
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}

let currentLeague = null;

function viewLeague(league) {
    currentLeague = league;
    document.getElementById('league-display').style.display = 'block';
    document.getElementById('league-title-display').textContent = league.title;
    document.getElementById('league-description-display').textContent = league.description;
}

async function deleteLeague() {
    if (!currentLeague) {
        alert('No league selected for deletion.');
        return;
    }

    try {
        const res = await fetch(`/api/leagues/${currentLeague.id}`, {
            method: 'DELETE',
            headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
        });
        if (!res.ok) {
            throw new Error('Failed to delete league.');
        }
        alert('League deleted successfully.');
        document.getElementById('league-display').style.display = 'none';
        await getLeagues();
    } catch (error) {
        alert(`Error: ${error.message}`);
    }
}

function setUploadButtonState(uploading, selector) {
    const uploadBtn = document.getElementById(selector);
    if (uploading) {
        uploadBtn.textContent = 'Uploading...';
        uploadBtn.disabled = true;
        return;
    }
    uploadBtn.textContent = 'Upload';
    uploadBtn.disabled = false;
}

async function uploadTournament(leagueID) {
    const excelSheets = document.getElementById('excel-sheets-selection').value;
    const categoryName = document.getElementById('tournament-category-name').value;
    const categoryDesc = document.getElementById('tournament-category-desc').value;
    const tournamentFile = document.getElementById('excel').files[0];
    if (!tournamentFile) return;

    const formData = new FormData();
    formData.append('excel', tournamentFile);
    formData.append('data', JSON.stringify({ excelSheets, categoryName, categoryDesc }));

    uploadBtnSelector = 'upload-excel-btn';
    setUploadButtonState(true, uploadBtnSelector);

    try {
        const res = await fetch(`/api/upload_tournament/${leagueID}`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
            body: formData,
        });
        if (!res.ok) {
            const data = await res.json();
            throw new Error(`Failed to upload excel. Error: ${data.error}`);
        }

        const bsc_response = await res.json();
        console.log('Tournament uploaded!');
        console.log(bsc_response);

        const bsc_output = document.getElementById('bsc-response');
        bsc_output.innerHTML = '';
        const contentItem = document.createElement('p');
        contentItem.textContent = bsc_response.bsc_reply;
        bsc_output.appendChild(contentItem);

        await getLeague(leagueID);
        await getLeageStandings(leagueID);
    } catch (error) {
        alert(`Error: ${error.message}`);
    }

    setUploadButtonState(false, uploadBtnSelector);
}


