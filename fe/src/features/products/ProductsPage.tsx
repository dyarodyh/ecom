import { useEffect, useState } from 'react'
import { getProducts } from './api'
import type { Product } from './types'
import CreateProductForm from './CreateProductForm'





type ProductsPageProps = {
  onAddToCart: (product: Product) => void
}

const rupiah = new Intl.NumberFormat('id-ID', {
  style: 'currency',
  currency: 'IDR',
  maximumFractionDigits: 0,
})


export default function ProductsPage({
  onAddToCart,
}: ProductsPageProps) {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  
function handleProductCreated(product: Product) {
  setProducts((currentProducts)=> [
    ...currentProducts,
    product,
  ])
}

  useEffect(() => {
    const controller = new AbortController()

    async function loadProducts() {
      try {
        const data = await getProducts(controller.signal)
        setProducts(data)
      } catch (err) {
        if (controller.signal.aborted) return

        setError(
          err instanceof Error
            ? err.message
            : 'Terjadi kesalahan saat mengambil produk.',
        )
      } finally {
        if (!controller.signal.aborted) {
          setLoading(false)
        }
      }
    }

    loadProducts()

    return () => {
      controller.abort()
    }
  }, [])

  if (loading) {
    return <p>Memuat produk...</p>
  }

  if (error) {
    return <p role="alert">{error}</p>
  }

  return (
    <section>
      <h1>Produk Toko Mini</h1>
      <CreateProductForm onCreated={handleProductCreated}/>

      {products.length === 0 ? (
        <p>Belum ada produk.</p>
      ) : (
        <ul>
          {products.map((product) => (
            <li key={product.id}>
              <h2>{product.name}</h2>
              <p>{rupiah.format(product.price)}</p>
              <p>Stok: {product.stock}</p>


              <button
              type='button'
              disabled={product.stock <= 0}
              onClick={() => onAddToCart(product)}
              >
                Tambah Ke Keranjang
              </button>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}