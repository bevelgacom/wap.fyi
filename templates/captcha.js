// Proof of Work Captcha System for Netscape 4+ compatibility
// wap.fyi captcha.js

var working = false;
var solution = 0;

// Simple hash function compatible with Netscape 4
function simpleHash(str) {
    var hash = 0;
    var i, char;
    if (str.length == 0) return hash;
    for (i = 0; i < str.length; i++) {
        char = str.charCodeAt(i);
        hash = ((hash << 5) - hash) + char;
        hash = hash & hash; // Convert to 32bit integer
    }
    // Make sure we get a positive number and add some variation
    hash = Math.abs(hash);
    if (hash == 0) hash = 1;
    return hash;
}

// Convert number to hex string with padding
function toHex(num) {
    var hex = num.toString(16);
    // Pad with zeros to ensure consistent length
    while (hex.length < 8) {
        hex = "0" + hex;
    }
    return hex;
}

// Check if hash has required number of trailing zeros
function hasTrailingZeros(hashNum, zeros) {
    var hexHash = toHex(hashNum);
    var trailingZeros = 0;
    for (var i = hexHash.length - 1; i >= 0 && trailingZeros < zeros; i--) {
        if (hexHash.charAt(i) == '0') {
            trailingZeros++;
        } else {
            break;
        }
    }
    return trailingZeros >= zeros;
}

// Start proof of work computation
function startProofOfWork(challengeFieldId, solutionFieldId, statusElementId, difficulty, maxIterations) {
    if (working) {
        alert("Already working! Please wait...");
        return;
    }
    
    var challengeField = document.getElementById(challengeFieldId);
    var challenge = challengeField ? challengeField.value : "";
    
    if (challenge == "") {
        alert("No challenge found! Please refresh the page.");
        return;
    }
    
    working = true;
    solution = 0;
    
    if (statusElementId) {
        var statusElement = document.getElementById(statusElementId);
        if (statusElement) {
            statusElement.innerHTML = "<span class='blink'>COMPUTING...</span> Please wait while your computer proves it's not a robot.";
            statusElement.className = "status working";
        }
    }
    
    // Update button visibility if function exists
    if (typeof updateButtonVisibility == 'function') {
        updateButtonVisibility();
    }
    
    // Use setTimeout to allow UI updates
    setTimeout(function() {
        doWork(challenge, solutionFieldId, statusElementId, difficulty || 4, maxIterations || 500);
    }, 100);
}

// Perform proof of work computation
function doWork(challenge, solutionFieldId, statusElementId, difficulty, maxIterations) {
    var iterations = 0;
    
    while (iterations < maxIterations && working) {
        var testString = challenge + solution;
        var hash = simpleHash(testString);
        
        if (hasTrailingZeros(hash, difficulty)) {
            // Found solution!
            working = false;
            
            // Set solution in hidden field
            var solutionField = document.getElementById(solutionFieldId);
            if (solutionField) {
                solutionField.value = solution;
            }
            
            if (statusElementId) {
                var statusElement = document.getElementById(statusElementId);
                if (statusElement) {
                    statusElement.innerHTML = "SUCCESS! You are verified!";
                    statusElement.className = "status success";
                }
            }
            
            // Call success callback if available
            if (typeof onProofOfWorkSuccess == 'function') {
                onProofOfWorkSuccess(challenge, solution);
            }
            
            return;
        }
        
        solution++;
        iterations++;
    }
    
    if (working) {
        // Update progress and continue
        if (statusElementId) {
            var statusElement = document.getElementById(statusElementId);
            if (statusElement) {
                var currentHash = simpleHash(challenge + solution);
                statusElement.innerHTML = "<span class='blink'>COMPUTING...</span> Tried " + solution + " possibilities... (Current hash: " + toHex(currentHash) + ")";
            }
        }
        setTimeout(function() {
            doWork(challenge, solutionFieldId, statusElementId, difficulty, maxIterations);
        }, 1);
    }
}

// Stop proof of work computation
function stopProofOfWork(statusElementId) {
    working = false;
    if (statusElementId) {
        var statusElement = document.getElementById(statusElementId);
        if (statusElement) {
            statusElement.innerHTML = "Computation stopped. Click Start Computation to try again.";
            statusElement.className = "status";
        }
    }
    
    // Update button visibility if function exists
    if (typeof updateButtonVisibility == 'function') {
        updateButtonVisibility();
    }
}

// Verify proof of work solution
function verifyProofOfWork(challengeFieldId, solutionFieldId, difficulty) {
    var challengeField = document.getElementById(challengeFieldId);
    var solutionField = document.getElementById(solutionFieldId);
    
    if (!challengeField || !solutionField) {
        return false;
    }
    
    var challenge = challengeField.value;
    var sol = parseInt(solutionField.value);
    
    if (challenge == "" || isNaN(sol)) {
        return false;
    }
    
    var testString = challenge + sol;
    var hash = simpleHash(testString);
    
    return hasTrailingZeros(hash, difficulty || 4);
}

// Reset proof of work (clear solution)
function resetProofOfWork(solutionFieldId, statusElementId) {
    working = false;
    solution = 0;
    
    var solutionField = document.getElementById(solutionFieldId);
    if (solutionField) {
        solutionField.value = "";
    }
    
    if (statusElementId) {
        var statusElement = document.getElementById(statusElementId);
        if (statusElement) {
            statusElement.innerHTML = "Click 'Start Computation' to begin verification.";
            statusElement.className = "status";
        }
    }
    
    // Update button visibility if function exists
    if (typeof updateButtonVisibility == 'function') {
        updateButtonVisibility();
    }
}
