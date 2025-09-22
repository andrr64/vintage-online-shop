# Checklist API Endpoints - E-commerce

## Manajemen Akun (/user)

- [ ] **Register** - `POST /user/register`
- [ ] **Login** - `POST /user/login`
- [ ] **Logout** - `POST /user/logout`
- [ ] **Manage Account** - `GET /user/account`
- [ ] **Update Profile** - `PUT /user/account/profile`  
  *Fields: username, fullname*
- [ ] **Change Password** - `PUT /user/account/password`

## Aktivitas Belanja (/shop)

- [ ] **Browse & Search Products** - `GET /products`
- [ ] **View Product Details** - `GET /products/{id}`
- [ ] **Manage Cart** - `POST/PUT/DELETE /cart`
- [ ] **Manage Wishlist** - `POST/DELETE /wishlist`
- [ ] **Checkout** - `POST /checkout`
- [ ] **View Order History** - `GET /orders`
- [ ] **Manage Addresses** - `GET/POST/PUT/DELETE /addresses`
