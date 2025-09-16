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
        dateCell.classList.add('transaction-date');

        const descriptionCell = document.createElement('td');
        descriptionCell.textContent = transaction.Description;
        descriptionCell.classList.add('transaction-description');

        const amountCell = document.createElement('td');
        amountCell.textContent = `$${transaction.Amount.toFixed(2)}`;
        amountCell.classList.add('transaction-amount');

        const categoryCell = document.createElement('td');
        categoryCell.textContent = transaction.Category;
        categoryCell.classList.add('transaction-category');

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
    if (!event.target.classList.contains('transaction-delete-btn')) {
        return;
    }
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
})

transactionsTableBody.addEventListener('click', (event) => {
    if (!event.target.classList.contains('transaction-edit-btn')) {
        return;
    }

    // Get the row that was clicked
    const row = event.target.closest('tr');

    // Get the data from the row
    const transactionId = event.target.dataset.id;
    const date = row.querySelector('.transaction-date').textContent;
    const description = row.querySelector('.transaction-description').textContent;
    const amount = row.querySelector('.transaction-amount').textContent.substring(1); // Remove the '$'
    const category = row.querySelector('.transaction-category').textContent;

    // Get references to the modal's form fields
    const editId = document.getElementById('edit-transaction-id');
    const editDate = document.getElementById('edit-transaction-date');
    const editDescription = document.getElementById('edit-transaction-description');
    const editAmount = document.getElementById('edit-transaction-amount');
    const editCategory = document.getElementById('edit-transaction-category');

    // Populate the form fields with the data
    editId.value = transactionId;
    editDate.value = date;
    editDescription.value = description;
    editAmount.value = parseFloat(amount);
    editCategory.value = category;

    // Show the modal dialog
    const editModalElement = document.getElementById('editTransactionModal');
    const editModal = new bootstrap.Modal(editModalElement);
    editModal.show();
});

const saveTransactionButton = document.getElementById('save-transaction-btn');

saveTransactionButton.addEventListener('click', (event) => {
    const transactionIdField = document.getElementById('edit-transaction-id');
    const transactionId = transactionIdField.value;

    console.log("Clicked Save for ID ", transactionId)
    
    fetch(`/api/transactions/${transactionId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(createTransactionJSONFromForm())
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to update the transaction.');
        }
        console.log("Transaction successfully updated.");
        return fetch('/api/transactions');
    })
    .then(response => response.json())
    .then(data => populateTransactionsGrid(data))
    .catch(error => {
        console.error('There was a problem with the update operation:', error);
    });
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

function createTransactionJSONFromForm() {
    const dateField = document.getElementById('edit-transaction-date');
    const descriptionField = document.getElementById('edit-transaction-description');
    const amountField = document.getElementById('edit-transaction-amount');
    const categoryField = document.getElementById('edit-transaction-category');

    return {
        date: dateField.value + 'T00:00:00Z',
        description: descriptionField.value,
        amount: parseFloat(amountField.value),
        category: categoryField.value
    };
}