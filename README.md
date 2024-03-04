# QR Menu Uygulaması

Bu QR Menu uygulaması, restoranların menülerini dijitalleştirmek ve müşterilere kolay bir şekilde erişim sağlamak için kullanılır. Bu uygulama, kullanıcıların kaydolmalarını, giriş yapmalarını ve menülerin görüntülemelerini sağlar. 

## Frontend
React Proje: [https://github.com/sencerarslan/qr-menu.git](https://github.com/sencerarslan/qr-menu.git)

## Demo

Demo: [https://qr-menu-wheat.vercel.app/](https://qr-menu-wheat.vercel.app/)

Demo Menu: [https://qr-menu-wheat.vercel.app/qr/65e5b0f2ae2dd5c28db6c381](https://qr-menu-wheat.vercel.app/qr/65e5b0f2ae2dd5c28db6c381)


## Kurulum

1. Bu projeyi klonlayın:

   ```bash
   git clone https://github.com/sencerarslan/go-app.git
   ```

2. Gerekli bağımlılıkları yükleyin:

   ```bash
   go mod tidy
   ```

3. MongoDB veritabanını çalıştırın ve bağlantı bilgilerini `databaseConnection.go` dosyasında güncelleyin.

4. Uygulamayı başlatmak için aşağıdaki komutu çalıştırın:

   ```bash
   go run main.go
   ```


## Teknolojiler

Bu proje aşağıdaki teknolojileri kullanır:

- Gin: Web framework olarak kullanılmıştır.
- MongoDB Driver: MongoDB veritabanı ile iletişim kurmak için kullanılmıştır.
- JWT-Go: JSON Web Token (JWT) oluşturmak ve doğrulamak için kullanılmıştır.
- Validator: Kullanıcı girdilerinin doğrulanması için kullanılmıştır.
- Godotenv: Çevresel değişkenlerin yüklenmesi için kullanılmıştır.