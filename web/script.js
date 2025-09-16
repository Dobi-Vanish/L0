document.addEventListener('DOMContentLoaded', function() {
    const orderIdInput = document.getElementById('orderId');
    const searchBtn = document.getElementById('searchBtn');
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');
    const result = document.getElementById('result');
    const orderDetails = document.getElementById('orderDetails');
    const jsonView = document.getElementById('jsonView');

    const apiBaseUrl = window.location.origin.replace(/:\d+/, ':8081');

    searchBtn.addEventListener('click', searchOrder);
    orderIdInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            searchOrder();
        }
    });

    function searchOrder() {
        const orderId = orderIdInput.value.trim();
        if (!orderId) {
            showError('Please enter an order ID');
            return;
        }

        resetUI();
        showLoading();

        fetch(`${apiBaseUrl}/order/${orderId}`)
            .then(response => {
                if (!response.ok) {
                    if (response.status === 404) {
                        throw new Error('Order not found');
                    } else {
                        throw new Error('Server error: ' + response.status);
                    }
                }
                return response.json();
            })
            .then(order => {
                hideLoading();
                displayOrder(order);
            })
            .catch(error => {
                hideLoading();
                showError(error.message);
            });
    }

    function displayOrder(order) {
        let html = `
            <div class="section">
                <h3 class="section-title">Order Information</h3>
                <p><strong>Order UID:</strong> ${order.order_uid}</p>
                <p><strong>Track Number:</strong> ${order.track_number}</p>
                <p><strong>Entry:</strong> ${order.entry}</p>
                <p><strong>Customer ID:</strong> ${order.customer_id}</p>
                <p><strong>Delivery Service:</strong> ${order.delivery_service}</p>
                <p><strong>Date Created:</strong> ${new Date(order.date_created).toLocaleString()}</p>
            </div>
            
            <div class="section">
                <h3 class="section-title">Delivery Information</h3>
                <p><strong>Name:</strong> ${order.delivery.name}</p>
                <p><strong>Phone:</strong> ${order.delivery.phone}</p>
                <p><strong>Zip:</strong> ${order.delivery.zip}</p>
                <p><strong>City:</strong> ${order.delivery.city}</p>
                <p><strong>Address:</strong> ${order.delivery.address}</p>
                <p><strong>Region:</strong> ${order.delivery.region}</p>
                <p><strong>Email:</strong> ${order.delivery.email}</p>
            </div>
            
            <div class="section">
                <h3 class="section-title">Payment Information</h3>
                <p><strong>Transaction:</strong> ${order.payment.transaction}</p>
                <p><strong>Currency:</strong> ${order.payment.currency}</p>
                <p><strong>Provider:</strong> ${order.payment.provider}</p>
                <p><strong>Amount:</strong> ${order.payment.amount}</p>
                <p><strong>Payment Date:</strong> ${new Date(order.payment.payment_dt * 1000).toLocaleString()}</p>
                <p><strong>Bank:</strong> ${order.payment.bank}</p>
                <p><strong>Delivery Cost:</strong> ${order.payment.delivery_cost}</p>
                <p><strong>Goods Total:</strong> ${order.payment.goods_total}</p>
                <p><strong>Custom Fee:</strong> ${order.payment.custom_fee}</p>
            </div>
        `;

        if (order.items && order.items.length > 0) {
            html += `
                <div class="section">
                    <h3 class="section-title">Items (${order.items.length})</h3>
                    ${order.items.map(item => `
                        <div style="margin-bottom: 15px; padding: 10px; border: 1px solid #eee; border-radius: 4px;">
                            <p><strong>Name:</strong> ${item.name}</p>
                            <p><strong>Brand:</strong> ${item.brand}</p>
                            <p><strong>Price:</strong> ${item.price}</p>
                            <p><strong>Total Price:</strong> ${item.total_price}</p>
                            <p><strong>Size:</strong> ${item.size}</p>
                            <p><strong>Status:</strong> ${item.status}</p>
                        </div>
                    `).join('')}
                </div>
            `;
        }

        orderDetails.innerHTML = html;
        jsonView.textContent = JSON.stringify(order, null, 2);
        result.classList.remove('hidden');
    }

    function showLoading() {
        loading.classList.remove('hidden');
    }

    function hideLoading() {
        loading.classList.add('hidden');
    }

    function showError(message) {
        error.textContent = message;
        error.classList.remove('hidden');
    }

    function resetUI() {
        error.classList.add('hidden');
        error.textContent = '';
        result.classList.add('hidden');
        orderDetails.innerHTML = '';
        jsonView.textContent = '';
    }
});