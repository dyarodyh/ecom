import ProductsPage from '../features/products/ProductsPage'
import CartPanel from '../features/cart/CartPanel'
import { useCart } from '../features/cart/useCart'
import './App.css'

export default function App() {

  const cart = useCart()
  return (
    <main>
      <ProductsPage onAddToCart={cart.addItem} />
      <hr />

      <CartPanel
        items={cart.items}
        totalPrice={cart.totalPrice}
        onRemoveItem={cart.removeItem}
      />
    </main>
  )
}