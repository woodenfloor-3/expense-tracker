package main

import (
    "html/template"
    "net/http"
    "os"
    "io"
    "strconv"
)

type Expense struct {
    Name     string
    Amount   float64
    Category string
    Image    string
}

var expenses []Expense

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/upload", uploadHandler)
    http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        r.ParseForm()

        // Convert amount from string to float64
        amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
        if err != nil {
            http.Error(w, "Invalid amount", http.StatusBadRequest)
            return
        }

        expense := Expense{
            Name:     r.FormValue("name"),
            Amount:   amount,
            Category: r.FormValue("category"),
            Image:    r.FormValue("image"),
        }
        expenses = append(expenses, expense)
    }

    tmpl := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
        <script src="https://unpkg.com/htmx.org@1.6.1"></script>
        <title>Expense Tracker</title>
    </head>
    <body class="bg-gray-100 p-6">
        <h1 class="text-2xl font-bold mb-4">Expense Tracker</h1>
        <form method="POST" enctype="multipart/form-data" class="mb-4">
            <div class="mb-2">
                <label>Name:</label>
                <input type="text" name="name" class="border rounded p-2 w-full" required>
            </div>
            <div class="mb-2">
                <label>Amount:</label>
                <input type="number" name="amount" class="border rounded p-2 w-full" required>
            </div>
            <div class="mb-2">
                <label>Category:</label>
                <select name="category" class="border rounded p-2 w-full">
                    <option value="Food">Food</option>
                    <option value="Transport">Transport</option>
                    <option value="Utilities">Utilities</option>
                </select>
            </div>
            <div class="mb-2">
                <label>Image:</label>
                <input type="file" name="image" class="border rounded p-2 w-full" required>
            </div>
            <button type="submit" class="bg-blue-500 text-white rounded p-2">Add Expense</button>
        </form>
        <table class="min-w-full bg-white border border-gray-300">
            <thead>
                <tr>
                    <th class="border px-4 py-2">Name</th>
                    <th class="border px-4 py-2">Amount</th>
                    <th class="border px-4 py-2">Category</th>
                    <th class="border px-4 py-2">Image</th>
                </tr>
            </thead>
            <tbody>
                {{range .}}
                <tr>
                    <td class="border px-4 py-2">{{.Name}}</td>
                    <td class="border px-4 py-2">{{.Amount}}</td>
                    <td class="border px-4 py-2">{{.Category}}</td>
                    <td class="border px-4 py-2"><img src="{{.Image}}" width="100"></td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </body>
    </html>
    `
    t := template.Must(template.New("webpage").Parse(tmpl))
    t.Execute(w, expenses)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(10 << 20) // limit your max input length!
    file, _, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Invalid image", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Create a temporary file
    tempFile, err := os.Create("uploads/" + "uploaded_image.jpg")
    if err != nil {
        http.Error(w, "Unable to save image", http.StatusInternalServerError)
        return
    }
    defer tempFile.Close()

    // Write the file to the temporary directory
    _, err = io.Copy(tempFile, file)
    if err != nil {
        http.Error(w, "Unable to save image", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}
