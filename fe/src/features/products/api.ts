import type {Product} from './types'


export type CreatedProductInput = {
    name: string
    price: number
    stock: number
}

export async function getProducts(
    signal?: AbortSignal,
): Promise<Product[]> {
    const response = await fetch ('/api/products',{
        signal,
    })

    if(!response.ok) {
        throw new Error (
            `Gagal megambil produk. Status HTTP: ${response.status}`,
        )
    }
    return response.json()
}

export async function createProduct(
    input: CreatedProductInput,
): Promise<Product> {
    const response = await fetch('/api/products', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(input),
    })

    if (!response.ok) {
        const message = (await response.text()).trim()

        throw new Error(
            message || `Gagal menambahkan produk. Status HTTP: ${response.status}`,
        )
    }
    return response.json()
}