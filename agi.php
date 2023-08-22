#!/usr/bin/php -q
<?php
require('phpagi.php');

$agi = new AGI();
$caller = "SOFTPHONE_B";
$pass = $agi->database_get("password",$caller);

$agi->answer();
$agi->exec("PLAYBACK", "welcome");

pwcheck($pass, $agi);

ivr:
$option = $agi->get_data("sekretarica", 20000, 1);
switch ($option['result']){
    case 1:
        $exten = $agi->get_data("vm-extension", 20000, 3);
        $agi->set_variable("tocall", $exten['result']);
        $agi->set_variable("redirect", "1");
        break;
    case 2:
        pwcheck($pass, $agi);
        $pw1 = $agi->get_data("vm-newpassword", 20000, 4);
        $pw2 = $agi->get_data("vm-reenterpassword", 20000, 4);
        if ($pw1 != $pw2 ){
            $agi->exec("PLAYBACK", "vm-mismatch");
            $agi->hangup();
        }
        $agi->database_put("password", $caller, $pw1['result']);
        $agi->exec("PLAYBACK", "vm-passchanged");
        break;
    case 3:
        goto ivr;
}

function pwcheck($pass, $agi){
    if(!$pass['result']) return;

    $result = $agi->get_data("vm-password", 9000, 4);
    if ($result['result'] != $pass['data']) {
        $agi->exec("PLAYBACK", "auth-incorrect");
        $agi->hangup();
    }
}
?>