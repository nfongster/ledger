const fetchButton = document.getElementById('fetch-budgets-btn');
const budgetsTableBody = document.getElementById('budgets-table-body');

fetchButton.addEventListener('click', () => {
    fetch('/api/budgets/status')
        .then(response => response.json())
        .then(data => {
            console.log(data);
            budgetsTableBody.innerHTML = '';

            data.forEach(budget => {
                const row = document.createElement('tr');

                const categoryCell = document.createElement('td');
                categoryCell.textContent = budget.Category;

                const periodCell = document.createElement('td');
                periodCell.textContent = budget.TimePeriod;

                const startDateCell = document.createElement('td');
                const startDateObject = new Date(budget.StartDate);
                startDateCell.textContent = startDateObject.toISOString().slice(1, 10);

                const endDateCell = document.createElement('td');
                const endDateObject = new Date(budget.EndDate);
                endDateCell.textContent = endDateObject.toISOString().slice(1, 10);

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

                row.appendChild(categoryCell);
                row.appendChild(periodCell);
                row.appendChild(startDateCell);
                row.appendChild(endDateCell);
                row.appendChild(budgetedAmountCell);
                row.appendChild(currentSpendingCell);
                row.appendChild(remainingCell);

                budgetsTableBody.appendChild(row);
            })
        })
        .catch(error => console.error('There was a problem with the fetch operation:', error));
});