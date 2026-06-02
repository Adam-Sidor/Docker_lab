<!DOCTYPE html>
<html lang="pl">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LEMP Stack Status</title>
</head>

<body>
    <div class="container">
        <h1>LEMP Stack Status</h1>

        <?php
        $host = 'mysql';
        $user = 'test_user';
        $pass = 'test_password';
        $db = 'test_db';

        try {
            $conn = new PDO("mysql:host=$host;dbname=$db", $user, $pass);
            $conn->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
            echo '<div class="status-box success">✔ Połączenie z bazą danych MySQL udane!</div>';
        } catch (PDOException $e) {
            echo '<div class="status-box error">❌ Błąd połączenia z bazą danych: ' . htmlspecialchars($e->getMessage()) . '</div>';
        }
        ?>

        <div class="info">
            <h3>Konfiguracja środowiska:</h3>
            <ul>
                <li><strong>PHP Version:</strong> <?php echo phpversion(); ?></li>
                <li><strong>Web Server:</strong> Nginx (Alpine)</li>
                <li><strong>Database Host:</strong> <?php echo $host; ?></li>
                <li><strong>Database Name:</strong> <?php echo $db; ?></li>
            </ul>
        </div>
    </div>
</body>

</html>