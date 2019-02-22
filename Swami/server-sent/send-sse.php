<?php
    // Setting up server
    header('Content-Type: text/event-stream');
    header('Cache-Control: no-cache');

    // Server Functionality 
        // Generates random number 0 to 1000 and sends it over
    $new_data = rand(0,1000);
    echo "data: New random number: $new_data";
    flush();

?>