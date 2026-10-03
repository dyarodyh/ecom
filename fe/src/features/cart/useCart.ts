import { useState } from 'react'
import type { Product } from '../products/types'
import type { CartItem } from './types'

export function useCart() {
  const [items, setItems] = useState<CartItem[]>([])

  function addItem(product: Product) {
    if (product.stock <= 0) return

    setItems((currentItems) => {
      const existingItem = currentItems.find(
        (item) => item.product.id === product.id,
      )

      if (existingItem) {
        if (existingItem.quantity >= product.stock) {
          return currentItems
        }

        return currentItems.map((item) =>
          item.product.id === product.id
            ? { ...item, quantity: item.quantity + 1 }
            : item,
        )
      }

      return [
        ...currentItems,
        {
          product,
          quantity: 1,
        },
      ]
    })
  }

  function removeItem(productId: number) {
    setItems((currentItems) =>
      currentItems.filter(
        (item) => item.product.id !== productId,
      ),
    )
  }

  const totalPrice = items.reduce(
    (total, item) =>
      total + item.product.price * item.quantity,
    0,
  )

  return {
    items,
    addItem,
    removeItem,
    totalPrice,
  }
}