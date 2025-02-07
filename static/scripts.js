updateTracksTable();

function createTrack() {
    const titleInput = document.getElementById("newTitle");

    fetch("/api/tracks/", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            data: {
                attributes: {
                    title: titleInput.value,
                },
            },
        }),
    })
        .then(updateTracksTable)
        .then((_) => (titleInput.value = ""));
}

function deleteTrack(element) {
    const row = element.parentNode.parentNode;
    const id = row.children[0].innerHTML;

    fetch("/api/tracks/" + id, {
        method: "DELETE",
    }).then(updateTracksTable);
}

function updateTracksTable() {
    fetch("/api/tracks/")
        .then((response) => response.json())
        .then((tracks) => {
            const tracksTable = document
                .getElementById("tracksTable")
                .getElementsByTagName("tbody")[0];
            tracksTable.innerHTML = "";

            tracks["data"].forEach((track) => {
                const newRow = tracksTable.insertRow();
                const id = newRow.insertCell();
                id.innerHTML = track["id"];

                const title = newRow.insertCell();
                title.innerHTML = track["attributes"]["title"];

                const removeBtn = newRow.insertCell();
                removeBtn.classList.add("action-cell");
                removeBtn.innerHTML = `<span class="delete-icon" onclick="deleteTrack(this)">✖</span>`;
            });
        });
}
