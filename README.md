# Go App

Bu Go uygulaması, kullanıcıların kayıt olmalarını, giriş yapmalarını ve kullanıcı verilerini yönetmelerini sağlayan basit bir web uygulamasıdır. Ayrıca, MongoDB veritabanı kullanılarak kullanıcı verileri depolanır.

## Kurulum

1. Bu projeyi klonlayın:

   ```bash
   git clone <https://github.com/sencerarslan/go-app.git>
   ```

2. Gerekli bağımlılıkları yükleyin:

   ```bash
   go mod tidy
   ```

3. MongoDB veritabanını çalıştırın ve bağlantı bilgilerini `database.go` dosyasında güncelleyin.

4. Uygulamayı başlatmak için aşağıdaki komutu çalıştırın:

   ```bash
   go run main.go
   ```

## Kullanım

- `POST /signup`: Yeni bir kullanıcı kaydı oluşturur. JSON formatında kullanıcı bilgilerini alır.

- `POST /login`: Kullanıcı girişi yapar. JSON formatında kullanıcı bilgilerini alır.

- `GET /users`: Tüm kullanıcıları listeler.

- `GET /users/:user_id`: Belirli bir kullanıcıyı alır.

## Gözden Geçirme

Bu projenin belirli işlevlerini kullanabilmek için aşağıdaki adımları izleyin:

1. `Signup` fonksiyonu, yeni bir kullanıcı kaydı oluşturur. Girilen kullanıcı bilgilerini doğrular ve veritabanında kontrol eder.

2. `Login` fonksiyonu, kullanıcı girişi yapar. Girilen bilgileri kontrol eder ve doğrulama işlemi yapar.

3. `GetUsers` fonksiyonu, tüm kullanıcıları listeler. İstenirse sayfalama yapılabilir.

4. `GetUser` fonksiyonu, belirli bir kullanıcıyı alır ve geri döndürür.

## Teknolojiler

Bu proje aşağıdaki teknolojileri kullanır:

- Gin: Web framework olarak kullanılmıştır.
- MongoDB Driver: MongoDB veritabanı ile iletişim kurmak için kullanılmıştır.
- JWT-Go: JSON Web Token (JWT) oluşturmak ve doğrulamak için kullanılmıştır.
- Validator: Kullanıcı girdilerinin doğrulanması için kullanılmıştır.
- Godotenv: Çevresel değişkenlerin yüklenmesi için kullanılmıştır.