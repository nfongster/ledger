document.addEventListener('DOMContentLoaded', () => {
    fetch('/api/budgets/status')
    .then(response => response.json())
    .then(data => populateBudgetsGrid(data))
    .catch(error => console.error('There was a problem with the fetch operation:', error));
});

const addBudgetButton = document.getElementById('add-budget-btn')

addBudgetButton.addEventListener('click', () => {
    fetch('/api/budgets', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(createBudgetJSON())
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP ${response.status} Error!  Message: ${response.statusText}`)
        }
        return response.json()
    })
    .then(data => {
        console.log(data);
        return fetch('/api/budgets');
    })
    .then(response => response.json())
    .then(data => populateBudgetsGrid(data))
    .catch(error => console.error('There was a problem with the fetch operation:', error));
});

const budgetsTableBody = document.getElementById('budgets-table-body');

function populateBudgetsGrid(budgetJsonArray) {
    console.log(budgetJsonArray);
    budgetsTableBody.innerHTML = '';

    budgetJsonArray.forEach(budget => {
        const row = document.createElement('tr');

        const categoryCell = document.createElement('td');
        categoryCell.textContent = budget.Category;
        categoryCell.classList.add('budget-category');

        const periodCell = document.createElement('td');
        periodCell.textContent = budget.TimePeriod;
        periodCell.classList.add('budget-period');

        const startDateCell = document.createElement('td');
        const startDateObject = new Date(budget.StartDate);
        startDateCell.textContent = startDateObject.toISOString().slice(0, 10);
        startDateCell.classList.add('budget-start-date');

        const endDateCell = document.createElement('td');
        const endDateObject = new Date(budget.EndDate);
        endDateCell.textContent = endDateObject.toISOString().slice(0, 10);
        endDateCell.classList.add('budget-end-date');

        const budgetedAmountCell = document.createElement('td');
        budgetedAmountCell.textContent = `$${budget.TargetAmount.toFixed(2)}`;
        budgetedAmountCell.classList.add('budget-target-amount');

        const currentSpendingCell = document.createElement('td');
        currentSpendingCell.textContent = `$${budget.CurrentSpent.toFixed(2)}`;

        const remainingCell = document.createElement('td');
        const remaining = budget.RemainingAmount.toFixed(2);
        

        if (budget.RemainingAmount < 0) {
            remainingCell.textContent = `-$${Math.abs(remaining)}`;
            remainingCell.style.color = 'red'
        } else if (budget.RemainingAmount > 0) {
            remainingCell.textContent = `$${remaining}`;
            remainingCell.style.color = 'green'
        }

        const editCell = document.createElement('td');
        const editButton = document.createElement('button');
        editButton.textContent = 'Edit';
        editButton.classList.add('btn', 'btn-warning', 'btn-sm', 'mb-3', 'budget-edit-btn');
        editButton.dataset.id = budget.ID;
        editCell.appendChild(editButton);

        const deleteCell = document.createElement('td');
        const deleteButton = document.createElement('button');
        deleteButton.textContent = 'Delete';
        deleteButton.classList.add('btn', 'btn-danger', 'btn-sm', 'mb-3', 'budget-delete-btn');
        deleteButton.dataset.id = budget.ID;
        deleteCell.appendChild(deleteButton);

        row.appendChild(categoryCell);
        row.appendChild(periodCell);
        row.appendChild(startDateCell);
        row.appendChild(endDateCell);
        row.appendChild(budgetedAmountCell);
        row.appendChild(currentSpendingCell);
        row.appendChild(remainingCell);
        row.appendChild(editCell);
        row.appendChild(deleteCell);

        budgetsTableBody.appendChild(row);
    })
}

budgetsTableBody.addEventListener('click', (event) => {
    if (!event.target.classList.contains('budget-delete-btn')) {
        return;
    }
    const budgetId = event.target.dataset.id;
    console.log("Clicked Delete for ID ", budgetId)
    
    fetch(`/api/budgets/${budgetId}`, {
        method: 'DELETE'
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to delete the budget.');
        }
        console.log("Budget successfully deleted.");
        return fetch('/api/budgets/status');
    })
    .then(response => response.json())
    .then(data => populateBudgetsGrid(data))
    .catch(error => {
        console.error('There was a problem with the delete operation:', error);
    });
})

budgetsTableBody.addEventListener('click', (event) => {
    if (!event.target.classList.contains('budget-edit-btn')) {
        return;
    }

    // Get the row that was clicked
    const row = event.target.closest('tr');

    // Get the data from the row
    const budgetId = event.target.dataset.id;
    const category = row.querySelector('.budget-category').textContent;
    const period = row.querySelector('.budget-period').textContent;
    const startDate = row.querySelector('.budget-start-date').textContent;
    const endDate = row.querySelector('.budget-end-date').textContent;
    const amount = row.querySelector('.budget-target-amount').textContent.substring(1); // Remove the '$'
    

    // Get references to the modal's form fields
    const editId = document.getElementById('edit-budget-id');
    const editCategory = document.getElementById('edit-budget-category');
    const editPeriod = document.getElementById('edit-budget-period');
    const editStartDate = document.getElementById('edit-budget-start-date');
    const editEndDate = document.getElementById('edit-budget-end-date');
    const editAmount = document.getElementById('edit-budget-amount');
    

    // Populate the form fields with the data
    editId.value = budgetId;
    editCategory.value = category;
    editPeriod.value = period;
    editStartDate.value = startDate;
    editEndDate.value = endDate;
    editAmount.value = parseFloat(amount);

    // Show the modal dialog
    const editModalElement = document.getElementById('editBudgetModal');
    const editModal = new bootstrap.Modal(editModalElement);
    editModal.show();
});

function createBudgetJSON() {
    const amountField = document.getElementById('budget-amount');
    const timePeriodField = document.getElementById('budget-period');
    const startDateField = document.getElementById('budget-start-date');
    const categoryField = document.getElementById('budget-category');
    
    return {
        target_amount: parseFloat(amountField.value),
        time_period: timePeriodField.value,
        start_date: startDateField.value + 'T00:00:00Z',
        category: categoryField.value
    };
}