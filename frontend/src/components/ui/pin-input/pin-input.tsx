import { useEffect, useRef, useState } from "react"

interface PinInputProps {
  length?: number
  value: string
  onChange: (value: string) => void
  onComplete?: (value: string) => void
}

export function PinInput({ length = 6, value, onChange, onComplete }: PinInputProps) {
  const [focusedIndex, setFocusedIndex] = useState<number>(0)
  const inputRefs = useRef<(HTMLInputElement | null)[]>([])

  useEffect(() => {
    if (inputRefs.current[focusedIndex]) {
      inputRefs.current[focusedIndex]?.focus()
    }
  }, [focusedIndex])

  useEffect(() => {
    if (value.length === length && onComplete) {
      onComplete(value)
    }
  }, [value, length, onComplete])

  const handleChange = (index: number, digit: string) => {
    if (!/^\d*$/.test(digit)) return

    const newValue = value.split("")
    newValue[index] = digit
    const updatedValue = newValue.join("").slice(0, length)

    onChange(updatedValue)

    if (digit && index < length - 1) {
      setFocusedIndex(index + 1)
    }
  }

  const handleKeyDown = (index: number, e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Backspace") {
      if (!value[index] && index > 0) {
        const newValue = value.split("")
        newValue[index - 1] = ""
        onChange(newValue.join(""))
        setFocusedIndex(index - 1)
      } else {
        const newValue = value.split("")
        newValue[index] = ""
        onChange(newValue.join(""))
      }
    } else if (e.key === "ArrowLeft" && index > 0) {
      setFocusedIndex(index - 1)
    } else if (e.key === "ArrowRight" && index < length - 1) {
      setFocusedIndex(index + 1)
    }
  }

  const handlePaste = (e: React.ClipboardEvent) => {
    e.preventDefault()
    const pastedData = e.clipboardData.getData("text").replace(/\D/g, "").slice(0, length)
    onChange(pastedData)
    setFocusedIndex(Math.min(pastedData.length, length - 1))
  }

  return (
    <div className="flex gap-2 items-center justify-center w-full">
      {Array.from({ length }).map((_, index) => {
        const hasValue = index < value.length
        const isActive = index === focusedIndex
        const isFilled = index < value.length - 1 || (index === value.length - 1 && hasValue)

        const baseClass =
          "w-12 h-14 rounded-xl flex items-center justify-center text-xl font-semibold outline-none transition-all text-center"
        const borderClass = isActive
          ? "border-2 border-primary-brand bg-card"
          : "border border-border"
        const bgClass = isFilled ? "bg-muted/30" : "bg-card"

        return (
          <div key={index} className="relative">
            <input
              ref={(el) => {
                inputRefs.current[index] = el
              }}
              type="text"
              inputMode="numeric"
              maxLength={1}
              value={value[index] || ""}
              onChange={(e) => handleChange(index, e.target.value)}
              onKeyDown={(e) => handleKeyDown(index, e)}
              onPaste={handlePaste}
              onFocus={() => setFocusedIndex(index)}
              className={`${baseClass} ${borderClass} ${bgClass}`}
              style={{
                caretColor: "transparent",
                color: "transparent",
              }}
            />
            {hasValue && (
              <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
                <span className="text-2xl font-semibold text-foreground">•</span>
              </div>
            )}
          </div>
        )
      })}
    </div>
  )
}
