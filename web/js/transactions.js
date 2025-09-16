document.addEventListener('DOMContentLoaded', () => {
    fetch('/api/transactions')
    .then(response => response.json())
    .then(data => {
        populateTransactionsGrid(data);
    })
    .catch(error => console.error('There was a problem with the fetch operation:', error));
});

const addTransactionButton = document.getElementById('add-transaction-btn');

addTransactionButton.addEventListener('click', () => {
    fetch('api/transactions', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(createTransactionJSON())
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP ${response.status} Error!  Message: ${response.statusText}`)
        }
        return response.json()
    })
    .then(data => {
        console.log(data);
        return fetch('/api/transactions');
    })
    .then(response => response.json())
    .then(data => populateTransactionsGrid(data))
    .catch(error => console.error('There was a problem with the fetch operation:', error));
});

const transactionsTableBody = document.getElementById('transactions-table-body');

function populateTransactionsGrid(transactionJsonArray) {
    console.log(transactionJsonArray); 
    transactionsTableBody.innerHTML = '';

    transactionJsonArray.forEach(transaction => {
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

        const editCell = document.createElement('td');
        const editButton = document.createElement('button');
        editButton.textContent = 'Edit';
        editButton.classList.add('btn', 'btn-warning', 'btn-sm', 'mb-3', 'transaction-edit-btn');
        editButton.dataset.id = transaction.ID;
        editCell.appendChild(editButton);

        const deleteCell = document.createElement('td');
        const deleteButton = document.createElement('button');
        deleteButton.textContent = 'Delete';
        deleteButton.classList.add('btn', 'btn-danger', 'btn-sm', 'mb-3', 'transaction-delete-btn');
        deleteButton.dataset.id = transaction.ID;
        deleteCell.appendChild(deleteButton);

        row.appendChild(dateCell);
        row.appendChild(descriptionCell);
        row.appendChild(amountCell);
        row.appendChild(categoryCell);
        row.appendChild(editCell);
        row.appendChild(deleteCell);

        transactionsTableBody.appendChild(row);
    })
}

transactionsTableBody.addEventListener('click', (event) => {
    if (event.target.classList.contains('transaction-delete-btn')) {
        const transactionId = event.target.dataset.id;
        console.log("Clicked Delete for ID ", transactionId)
        
        fetch(`/api/transactions/${transactionId}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to delete the transaction.');
            }
            console.log("Transaction successfully deleted.");
            return fetch('/api/transactions');
        })
        .then(response => response.json())
        .then(data => populateTransactionsGrid(data))
        .catch(error => {
            console.error('There was a problem with the delete operation:', error);
        });
    }
})

function createTransactionJSON() {
    const dateField = document.getElementById('transaction-date');
    const descriptionField = document.getElementById('transaction-description');
    const amountField = document.getElementById('transaction-amount');
    const categoryField = document.getElementById('transaction-category');

    return {
        date: dateField.value + 'T00:00:00Z',
        description: descriptionField.value,
        amount: parseFloat(amountField.value),
        category: categoryField.value
    };
}