import type { CartItem } from './types'

type CartPanelProps = {
    items: CartItem[]
    totalPrice: number
    onRemoveItem: (productId: number) => void
}

export default function CartPanel({
    items,
    totalPrice,
    onRemoveItem,
}: CartPanelProps) {
    return (
        <section>
            <h2>Keranjang Belanja</h2>

            {items.length === 0 ? (
                <p>Keranjang Masih Kosong.</p>
            ) : (
                <>
                    <ul>
                        {items.map((item) => (
                            <li key={item.product.name}>
                                <h3>{item.product.name}</h3>

                                <p>
                                    Jumlah: {item.quantity}
                                    {'/'}
                                    Stock: {item.product.stock}
                                </p>

                                <p>
                                    Subtotal: Rp
                                    {(
                                        item.product.price * item.quantity
                                    ).toLocaleString('id-ID')}
                                </p>

                                <button 
                                type='button'
                                onClick={()=>
                                    onRemoveItem(item.product.id)
                                }>
                                    Hapus
                                </button>
                            </li>
                        ))}
                    </ul>

                    <p>
                        <strong>
                            Total: Rp{totalPrice.toLocaleString('id-ID')}
                        </strong>
                    </p>
                </>
            )}
        </section>
    )
}

