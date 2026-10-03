export async function getSavedTeam () {
    const emptyTeam = [];
    try {
        const resp = await fetch("/api/getTeamNames", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
        });

        const data = await resp.json();
        if (!resp.ok) {
            console.log(`error: ${data.error}`);
            return emptyTeam;
        }

        if (data?.team == null || data?.team == undefined) {
            console.log(`Team is ${data?.team}. Returning empty team.`);
            return emptyTeam;
        }
        else {
            console.log(`Successfully retrieved team from server: ${data.team}`);
            return  data.team;
        }

    } catch (error) {
        alert(`Error: ${error.message}`);
        return emptyTeam;
    }
}    