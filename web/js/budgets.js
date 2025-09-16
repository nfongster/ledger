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

        const periodCell = document.createElement('td');
        periodCell.textContent = budget.TimePeriod;

        const startDateCell = document.createElement('td');
        const startDateObject = new Date(budget.StartDate);
        startDateCell.textContent = startDateObject.toISOString().slice(0, 10);

        const endDateCell = document.createElement('td');
        const endDateObject = new Date(budget.EndDate);
        endDateCell.textContent = endDateObject.toISOString().slice(0, 10);

        const budgetedAmountCell = document.createElement('td');
        budgetedAmountCell.textContent = `$${budget.TargetAmount.toFixed(2)}`;

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
    if (event.target.classList.contains('budget-delete-btn')) {
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
    }
})

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