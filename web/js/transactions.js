const fetchButton = document.getElementById('fetch-transactions-btn');
const transactionsTableBody = document.getElementById('transactions-table-body');
const deleteButtonsContainer = document.getElementById('delete-buttons-container');

fetchButton.addEventListener('click', () => {
    fetch('/api/transactions')
        .then(response => response.json())
        .then(data => {
            console.log(data); 
            transactionsTableBody.innerHTML = '';
            deleteButtonsContainer.innerHTML = '';

            data.forEach(transaction => {
                const row = document.createElement('tr');

                const dateCell = document.createElement('td');
                const dateObject = new Date(transaction.Date);
                dateCell.textContent = dateObject.toISOString().slice(0, 10);

                const descriptionCell = document.createElement('td');
                descriptionCell.textContent = transaction.Description;

                const amountCell = document.createElement('td');
                amountCell.textContent = `$${transaction.Amount.toFixed(2)}`;

                const categoryCell = document.createElement('td');
                categoryCell.textContent = transaction.Category;

                const deleteCell = document.createElement('td');

                const deleteButton = document.createElement('button');
                deleteButton.textContent = 'Delete';
                deleteButton.classList.add('btn', 'btn-danger', 'btn-sm', 'mb-3');
                deleteButton.dataset.id = transaction.id;

                deleteCell.appendChild(deleteButton);

                row.appendChild(dateCell);
                row.appendChild(descriptionCell);
                row.appendChild(amountCell);
                row.appendChild(categoryCell);
                row.appendChild(deleteCell);

                transactionsTableBody.appendChild(row);
            })
        })
        .catch(error => console.error('There was a problem with the fetch operation:', error));
});