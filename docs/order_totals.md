# OpenB2B Order & Invoice Calculations

This document explains how **order totals, item discounts, order-level discounts, taxes, and invoice totals** are calculated in OpenB2B.

---

## 1️⃣ Definitions

| Field                  | Description                                                                      |
| ---------------------- | -------------------------------------------------------------------------------- |
| `UnitPrice`            | Price of a single product variant.                                               |
| `Quantity`             | Number of units ordered for the item.                                            |
| `ItemDiscount`         | Discount applied to an individual order item. Can be percentage or fixed amount. |
| `AppliedOrderDiscount` | Portion of the **order-level discount** distributed to this item proportionally. |
| `TaxRate`              | Tax percentage for the item.                                                     |
| `Subtotal`             | Sum of all items’ `UnitPrice * Quantity` before discounts.                       |
| `ItemDiscountTotal`    | Sum of all per-item discounts.                                                   |
| `AppliedDiscount`      | Actual order-level discount applied after checking limits.                       |
| `DiscountTotal`        | `ItemDiscountTotal + AppliedDiscount`                                            |
| `TaxAmount`            | Tax calculated per item after all discounts.                                     |
| `Total`                | Final total per item including tax.                                              |
| `TaxTotal`             | Sum of all `TaxAmount` for the order.                                            |
| `TotalAmount`          | Final total for the order including tax.                                         |

---

## 2️⃣ Calculation Steps

### Step 1: Calculate Item Discount

For each order item:

```text
Subtotal = UnitPrice * Quantity

If ItemDiscount.Type == "percentage":
    AppliedDiscount = Subtotal * (ItemDiscount.Amount / 100)
Else if ItemDiscount.Type == "fixed":
    AppliedDiscount = ItemDiscount.Amount
Else:
    AppliedDiscount = 0

AppliedDiscount = min(AppliedDiscount, Subtotal)
```

---

### Step 2: Calculate Subtotal & ItemDiscountTotal

```text
Subtotal = sum(UnitPrice * Quantity for all items)
ItemDiscountTotal = sum(AppliedDiscount for all items)
```

---

### Step 3: Calculate Order-Level Discount

```text
If Order.Discount.Type == "percentage":
    AppliedDiscount = Subtotal * (Order.Discount.Amount / 100)
Else if Order.Discount.Type == "fixed":
    AppliedDiscount = Order.Discount.Amount
Else:
    AppliedDiscount = 0

MaxAllowedDiscount = Subtotal - ItemDiscountTotal
AppliedDiscount = min(AppliedDiscount, MaxAllowedDiscount)
```

---

### Step 4: Distribute Order Discount to Items

Each item receives a **proportional share** of the order-level discount:

```text
For each item:
    proportion = (UnitPrice * Quantity) / Subtotal
    AppliedOrderDiscount = proportion * AppliedDiscount

Adjust last item to absorb rounding differences to ensure:
sum(AppliedOrderDiscount for all items) == AppliedDiscount
```

---

### Step 5: Calculate Tax and Line Total

For each item:

```text
TaxableAmount = (UnitPrice * Quantity) - AppliedDiscount - AppliedOrderDiscount
If TaxableAmount < 0: TaxableAmount = 0

TaxAmount = round(TaxableAmount * TaxRate, 2)
LineTotal = round(TaxableAmount + TaxAmount, 2)
```

* `AppliedDiscount` = item-level discount
* `AppliedOrderDiscount` = proportional order-level discount

---

### Step 6: Aggregate Totals for the Order

```text
TaxTotal = sum(TaxAmount for all items)
Total = sum(LineTotal for all items)
DiscountTotal = ItemDiscountTotal + AppliedDiscount
```

---

## 3️⃣ Example Calculation

Order:

| Item | UnitPrice | Quantity | ItemDiscount | TaxRate |
| ---- | --------- | -------- | ------------ | ------- |
| 1    | 100       | 2        | 10%          | 10%     |
| 2    | 50        | 1        | \$5          | 5%      |

Order-level discount: \$20 (fixed)

**Step 1: Item Discounts**

* Item 1: 100\*2 \* 10% = 20
* Item 2: \$5

**Step 2: Subtotal & ItemDiscountTotal**

* Subtotal = 200 + 50 = 250
* ItemDiscountTotal = 20 + 5 = 25

**Step 3: Order-Level Discount**

* Fixed \$20 → AppliedDiscount = 20

**Step 4: Distribute Order Discount**

* Item 1: 200/250 \* 20 = 16
* Item 2: 50/250 \* 20 = 4

**Step 5: Tax & Line Totals**

* Item 1 taxable = 200 - 20 - 16 = 164 → Tax = 16.4 → Total = 180.4
* Item 2 taxable = 50 - 5 - 4 = 41 → Tax = 2.05 → Total = 43.05

**Step 6: Aggregate**

* TaxTotal = 16.4 + 2.05 = 18.45
* Total = 180.4 + 43.05 = 223.45
* DiscountTotal = 25 + 20 = 45

---

## 4️⃣ Notes

* **All rounding** is done to 2 decimal places (cents).
* **Order-level discount is capped** at `Subtotal - ItemDiscountTotal`.
* **Tax is calculated after all discounts**, which is standard practice.
* Proportional distribution ensures the sum of applied order discounts matches the order-level discount exactly.

---

This document can be included in your repo as **`docs/order_totals.md`** for reference.

---

If you want, I can also make a **diagram showing the flow** from **unit price → item discount → order discount → tax → totals**, which helps visualize the calculation.

Do you want me to do that?
