/**
 *
 * appCtrl
 *
 */

angular
    .module('homer')
    .controller('appCtrl', appCtrl)
    .controller('dashCtrl', dashCtrl);

function appCtrl($http, $scope) {
}

function dashCtrl($http, $scope, growl) {
    $scope.info             = {};
    $scope.spam_text        = "";
    $scope.pred             = {};
    $scope.pred.havepred    = false;
    $scope.pypred           = {};
    $scope.pypred.havepred  = false;

    $http({
        method: 'GET',
        url: '/getModelInfo'
      }).then(function successCallback(response) {
          console.log('success response: ' + JSON.stringify(response));
          $scope.info = response.data.content;
        }, function errorCallback(response) {
            console.log('error response: ' + JSON.stringify(response));
        });

    $scope.getPred = function() {
        $scope.pred.havepred    = false;
        $scope.pypred.havepred  = false;

        if ($scope.spam_text.length == 0 || $scope.spam_text == undefined) {
            growl.warning("no text entered, please enter submit again", {ttl: 2500});
            return 
        }

        $http({
            method: 'POST',
            data: {"text": $scope.spam_text},
            url: '/getModelPred'
          }).then(function successCallback(response) {
              console.log('success response: ' + JSON.stringify(response));
              $scope.pred = response.data.content;
              growl.success("Predictions Updated", {ttl: 1000});
            }, function errorCallback(response) {
                console.log('error response: ' + JSON.stringify(response));
                growl.warning(response.data.msg, {ttl: 2500});
            });

        $http({
            method: 'POST',
            data: {"text": $scope.spam_text},
            url: '/getPyModelPred'
            }).then(function successCallback(rsp) {
                console.log('success response: ' + JSON.stringify(rsp));
                $scope.pypred = rsp.data.content;
            }, function errorCallback(rsp) {
                console.log('error response: ' + JSON.stringify(rsp));
                growl.warning(rsp.data.msg, {ttl: 2500});
            });
    };

    $scope.fillTestMsg = function(val) {
        $scope.spam_text        = "";
        $scope.pred.havepred    = false;
        $scope.pypred.havepred  = false;

        if (val == "clear") {
            return; 
        }

        $http({
            method: 'GET',
            url: '/getRandMsg/' + val
          }).then(function successCallback(response) {
              console.log('success response: ' + JSON.stringify(response));
              $scope.spam_text = response.data.content;
              growl.success("Random " + val + " message", {ttl: 1000});
            }, function errorCallback(response) {
                console.log('error response: ' + JSON.stringify(response));
                growl.warning("an error occurred", {ttl: 1000});
            });
    };
}
