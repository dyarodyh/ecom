import { useState } from 'react'
import type { FormEvent } from 'react'
import { createProduct } from './api'
import type { Product } from './types'

type CreateProductFormProps = {
    onCreated: (product: Product) => void
}

export default function CreateProductForm({
    onCreated,
}: CreateProductFormProps) {
    const [isSaving, setIsSaving] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [success, setSuccess] = useState<string | null>(null)

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()

        if (isSaving) return

        setError(null)
        setSuccess(null)

        const form = event.currentTarget
        const fields = new FormData(form)

        const name = String(fields.get('name') ?? '').trim()
        const priceText = String(fields.get('price') ?? '').trim()
        const stockText = String(fields.get('stock') ?? '').trim()

        const price = Number(priceText)
        const stock = Number(stockText)

        if (name === '') {
            setError('Nama produk wajib diisi.')
            return
        }

        if (
            priceText === '' ||
            !Number.isSafeInteger(price) ||
            price < 0
        ) {
            setError('Harga harus berupa bilangan bulat 0 atau lebih.')
            return
        }

        if (
            stockText === '' ||
            !Number.isSafeInteger(stock) ||
            stock < 0 ||
            stock > 2_147_483_647
        ) {
            setError('Stok harus berupa bilangan bulat antara 0 dan 2147483647.')
            return
        }

        setIsSaving(true)

        try {
            const product = await createProduct({
                name,
                price,
                stock,
            })

            onCreated(product)
            form.reset()
            setSuccess(`Produk "${product.name}" berhasil ditambahkan.`)
        } catch (err) {
            setError(
                err instanceof Error
                    ? err.message
                    : 'Terjadi kesalahan saat menambahkan produk.',
            )
        } finally {
            setIsSaving(false)
        }
    }

    return (
        <form onSubmit={handleSubmit}>
            <fieldset disabled={isSaving}>
                <legend>Tambah produk baru</legend>

                <p>
                    <label>
                        Nama produk
                        <br />
                        <input name="name" type="text" required />
                    </label>
                </p>

                <p>
                    <label>
                        Harga (Rp)
                        <br />
                        <input
                            name="price"
                            type="number"
                            min="0"
                            max={Number.MAX_SAFE_INTEGER}
                            step="1"
                            required
                        />
                    </label>
                </p>

                <p>
                    <label>
                        Stok
                        <br />
                        <input
                            name="stock"
                            type="number"
                            min="0"
                            max="2147483647"
                            step="1"
                            defaultValue="0"
                            required
                        />
                    </label>
                </p>

                <button type="submit">
                    {isSaving ? 'Menyimpan...' : 'Simpan produk'}
                </button>
            </fieldset>

            {error && <p role="alert">{error}</p>}
            {success && <p role="status">{success}</p>}
        </form>
    )
}