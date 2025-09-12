const fetchButton = document.getElementById('fetch-transactions-btn');

fetchButton.addEventListener('click', () => {
    fetch('/api/transactions')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log(data); 
            // This is where you will add logic to populate your table
        })
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
        });
});

