document.addEventListener("DOMContentLoaded", function () {
    function fetchExpressions() {
        fetch("http://127.0.0.1:8080/api/v1/expressions")
            .then(response => response.json())
            .then(data => {
                let tableBody = document.getElementById("expressions-table");
                tableBody.innerHTML = "";

                data.expressions.forEach(expr => {
                    let row = `<tr>
                        <td>${expr.id}</td>
                        <td>${expr.status}</td>
                        <td>${expr.result || "—"}</td>
                    </tr>`;
                    tableBody.innerHTML += row;
                });
            })
            .catch(error => console.error("Error fetching expressions:", error));
    }

    setInterval(fetchExpressions, 1000);
    fetchExpressions();
});
